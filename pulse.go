package main

import (
	"github.com/AlpineMarmot/pulse/config"
	"github.com/AlpineMarmot/pulse/database"
	"github.com/AlpineMarmot/pulse/logger"
	"github.com/AlpineMarmot/pulse/util"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/redisstorage"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"runtime"
	"time"
)

type Pulse struct {
	entryPoint   string
	config       config.Config
	configFile   string
	colly        *colly.Collector
	collyStorage struct {
		redisStorage redisstorage.Storage
	}
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
		logger.Println("Loading default configuration ... ")
		p.config = config.NewConfigFromFile(util.GetCurrentPath() + "/default.yml")
	} else if util.FileExists(configFile) == true {
		logger.Println("Loading configuration file", configFile)
		p.config = config.NewConfigFromFile(configFile)
	} else if util.FileExists(configFile) == true {
		logger.Println(configFile)
		logger.Printf("Unable to load configuration file : %s", configFile)
		return
	}
	p.applyConfigToColly()
}

func (p *Pulse) CreateSessionID(db database.MongoDb) interface{} {
	coll := db.Collection("sessions")
	res, _ := coll.InsertOne(db.GetQueryContext(), bson.M{
		"dt_created": time.Now(),
	})
	return res.InsertedID
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
		logger.Println("Please, specify an url as entry point ...")
		return
	}

	logger.Println("Starting ...")

	p.colly.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
		time.Sleep(p.config.Crawler.SleepBetweenRequest * time.Millisecond)
	})

	for _, middleware := range p.requestMiddlewares {
		p.colly.OnRequest(middleware)
	}

	p.colly.OnResponse(func(r *colly.Response) {
		if p.maxUrlVisited > 0 && p.maxUrlVisited <= p.stats.urlVisited {
			logger.Printf("\nLimit of %d URLs visited reached!\n", p.maxUrlVisited)
			p.printStats()
			os.Exit(0)
		}
		p.stats.urlVisited++
		logger.Println("Visited link #", p.stats.urlVisited, r.Ctx.Get("url"))
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
	logger.Println("URL visited:", p.stats.urlVisited)
	logger.Println("URL skipped:", p.stats.urlSkipped)
	logger.Println("URL missing or partial:", p.stats.urlMissing)
	logger.Println("URL blocked by robots.txt:", p.stats.urlBlockedByRobotsTxt)
	now := time.Now()
	logger.Println(now.Sub(p.stats.startTime).Seconds(), "sec(s)")
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

	limitParallelism := p.config.Crawler.Limit.Parallelism
	if limitParallelism == 0 {
		limitParallelism = runtime.NumCPU()
	}

	err := p.colly.Limit(&colly.LimitRule{
		DomainGlob:   p.config.Crawler.Limit.DomainGlob,
		DomainRegexp: p.config.Crawler.Limit.DomainRegexp,
		Delay:        p.config.Crawler.Limit.Delay,
		RandomDelay:  p.config.Crawler.Limit.RandomDelay,
		Parallelism:  limitParallelism,
	})
	util.CheckError(err, "Applying Limit config to colly")

	p.colly.MaxBodySize = p.config.Crawler.MaxBodySize
	p.colly.MaxDepth = p.config.Crawler.MaxDepth

	p.setStorage()

	if p.config.Crawler.RandomUserAgents == true {
		extensions.RandomUserAgent(p.colly)
	}
}

func (p *Pulse) setStorage() {
	if len(p.config.Crawler.Storage.Redis.Address) > 0 {
		p.collyStorage.redisStorage = redisstorage.Storage{
			Address:  p.config.Crawler.Storage.Redis.Address,
			Password: p.config.Crawler.Storage.Redis.Password,
			DB:       p.config.Crawler.Storage.Redis.Db,
			Prefix:   p.config.Crawler.Storage.Redis.Prefix,
		}
		err := p.colly.SetStorage(&p.collyStorage.redisStorage)
		util.CheckError(err, "Setting redis as internal storage")

		if p.config.Crawler.Storage.ClearOnStart == true {
			err := p.collyStorage.redisStorage.Clear()
			util.CheckError(err, "Clearing redis storage")
		}
	}
}

func (p *Pulse) CloseStorage() {
	if len(p.config.Crawler.Storage.Redis.Address) > 0 {
		_ = p.collyStorage.redisStorage.Client.Close()
	}
}
