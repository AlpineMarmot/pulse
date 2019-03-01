package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"pulse/config"
	"pulse/util"
	"runtime"
	"strconv"
	"time"
)

type Pulse struct {
	entryPoint          string
	config              config.Config
	configFile          string
	colly               *colly.Collector
	requestMiddlewares  []colly.RequestCallback
	responseMiddlewares []colly.ResponseCallback
	htmlMiddlewares     map[string]colly.HTMLCallback
	maxUrlVisited       int64
	stats               Stats
}

type Stats struct {
	urlVisited            int64
	urlSkipped            int64
	urlBlockedByRobotsTxt int64
	urlMissing            int64
	startTime             time.Time
}

type Events struct {
}

func NewPulse() Pulse {
	return Pulse{
		colly:           colly.NewCollector(),
		htmlMiddlewares: make(map[string]colly.HTMLCallback),
		maxUrlVisited:   0,
		stats: Stats{
			urlVisited:            0,
			urlSkipped:            0,
			urlBlockedByRobotsTxt: 0,
			urlMissing:            0,
			startTime:             time.Now(),
		},
	}
}

func (p *Pulse) SetEntryPoint(entryPoint string) {
	p.entryPoint = entryPoint
}

func (p *Pulse) LoadConfigFile(configFile string) {
	if len(configFile) == 0 {
		fmt.Println("Loading default configuration ... ")
		p.config = config.NewConfigFromString(config.GetDefaultConfig())
	} else if config.ConfigFileExists(configFile) == true {
		fmt.Println("Loading configuration file", configFile)
		p.config = config.NewConfigFromFile(configFile)
	} else if config.ConfigFileExists(configFile) == true {
		fmt.Println(configFile)
		fmt.Printf("Unable to load configuration file : %s", configFile)
		return
	}

	p.applyConfigToColly()
}

func (p *Pulse) OnRequest(middleware colly.RequestCallback) {
	p.requestMiddlewares = append(p.requestMiddlewares, middleware)
}

func (p *Pulse) OnResponse(middleware colly.ResponseCallback) {
	p.responseMiddlewares = append(p.responseMiddlewares, middleware)
}

func (p *Pulse) OnHTML(goquerySelector string, middleware colly.HTMLCallback) {
	p.htmlMiddlewares[goquerySelector] = middleware
}

func (p *Pulse) Start() {
	if len(p.entryPoint) < 1 {
		fmt.Println("Please, specify an url as entry point ...")
		return
	}

	fmt.Println("Starting ...")

	p.colly.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		time.Sleep(p.config.Crawler.SleepBetweenRequest * time.Millisecond)
	})

	for _, middleware := range p.requestMiddlewares {
		p.colly.OnRequest(middleware)
	}

	p.colly.OnResponse(func(r *colly.Response) {
		if p.maxUrlVisited > 0 && p.maxUrlVisited <= p.stats.urlVisited {
			fmt.Printf("\nLimit of %d URLs visited reached!\n", p.maxUrlVisited)
			p.printStats()
			os.Exit(0)
		}
		p.stats.urlVisited++
		fmt.Println("\nVisited", p.stats.urlVisited, r.Ctx.Get("url"))
	})

	for _, middleware := range p.responseMiddlewares {
		p.colly.OnResponse(middleware)
	}

	p.colly.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		err := e.Request.Visit(e.Request.AbsoluteURL(url))
		if err != nil {
			if err.Error() == colly.ErrAlreadyVisited.Error() ||
				err.Error() == colly.ErrMissingURL.Error() ||
				err.Error() == colly.ErrRobotsTxtBlocked.Error() {
				p.statsFromError(err)
				return
			}
		}
		util.CheckError(err, "Jumping on a link "+url)
	})

	for selector, middleware := range p.htmlMiddlewares {
		p.colly.OnHTML(selector, middleware)
	}

	err := p.colly.Visit(p.entryPoint)
	util.CheckError(err, "Starting pulse")
	p.colly.Wait()

	p.printStats()
}

func (p *Pulse) statsFromError(err error) {
	switch err.Error() {
	case colly.ErrAlreadyVisited.Error():
		p.stats.urlSkipped++
	case colly.ErrMissingURL.Error():
		p.stats.urlMissing++
	case colly.ErrRobotsTxtBlocked.Error():
		p.stats.urlBlockedByRobotsTxt++
	}
}

func (p *Pulse) printStats() {
	fmt.Println("URL visited:", p.stats.urlVisited)
	fmt.Println("URL skipped:", p.stats.urlSkipped)
	fmt.Println("URL missing or partial:", p.stats.urlMissing)
	fmt.Println("URL blocked by robots.txt:", p.stats.urlBlockedByRobotsTxt)
	now := time.Now()
	fmt.Println(now.Sub(p.stats.startTime).Seconds(), "sec(s)")
}

func (p *Pulse) applyConfigToColly() {
	p.maxUrlVisited = p.config.Crawler.MaxUrlVisited
	p.colly.AllowedDomains = p.config.Crawler.AllowedDomains
	p.colly.AllowURLRevisit = p.config.Crawler.AllowURLRevisit
	p.colly.Async = p.config.Crawler.Async
	p.colly.CheckHead = p.config.Crawler.CheckHead
	p.colly.DetectCharset = p.config.Crawler.DetectCharset
	p.colly.DisallowedDomains = p.config.Crawler.DisallowedDomains
	p.colly.IgnoreRobotsTxt = p.config.Crawler.IgnoreRobotsTxt

	realLimitParallelism := 1
	limitParallelism := p.config.Crawler.Limit.Parallelism
	if limitParallelism == "auto" {
		realLimitParallelism = runtime.NumCPU()
	} else {
		x := p.config.Crawler.Limit.Parallelism.(string)
		realLimitParallelism, _ = strconv.Atoi(x) // Alt. non panicking version
	}

	fmt.Println(realLimitParallelism)
	err := p.colly.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: realLimitParallelism,
		//RandomDelay: 15 * time.Second,
	})
	util.CheckError(err, "Applying Limit config to colly")

	p.colly.MaxDepth = p.config.Crawler.MaxDepth

	// need work
	//if len(p.config.Crawler.Proxy) > 0 {
	//	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	p.colly.SetProxyFunc(rp)
	//}
}
