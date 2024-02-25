package persistence

import (
	"gorm.io/gorm"
)

type PagingOptions struct {
	Offset int
	Limit  int
}

type PagingConfigurator func(options *PagingOptions)

// Paginate returns a function that can be used to paginate a *gorm.DB object.
//
// It takes a variable number of PagingConfigurator functions as input and returns a function that accepts a *gorm.DB object
// and returns a modified *gorm.DB object with pagination applied.
//
// The PagingConfigurator functions are used to configure the pagination options such as offset and limit.
// The function applies the configured options to the *gorm.DB object and returns the modified object.
//
// The default pagination options are set to Offset: 0 and Limit: 20.
// The function loops through the provided configurators and calls each one with the config options.
// It then checks if the configured limit is greater than 50 and sets it to 50 if so.
// It also checks if the configured offset is less than 0 and sets it to 0 if so.
//
// The function returns the modified *gorm.DB object with the applied pagination options.
func Paginate(configurators ...PagingConfigurator) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		config := &PagingOptions{
			Offset: 0,
			Limit:  50,
		}

		// Apply the configurators
		for _, configurator := range configurators {
			configurator(config)
		}

		// Force sane limits
		switch {
		case config.Limit > 50:
			config.Limit = 50

		case config.Offset < 0:
			config.Offset = 0
		}

		return db.Offset(config.Offset).Limit(config.Limit)
	}
}
