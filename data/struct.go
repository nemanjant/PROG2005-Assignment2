package data

type RequestOpenMeteo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2M float64 `json:"temperature_2m"`
		Precipitation float64 `json:"precipitation"`
	} `json:"current"`
}

type Registration struct {
	Id       string `json:"id"`
	Country  string `json:"country"`
	ISOcode  string `json:"isoCode"`
	Features struct {
		Temperature      bool     `json:"temperature"`
		Precipitation    bool     `json:"precipitation"`
		Capital          bool     `json:"capital"`
		Coordinates      bool     `json:"coordinates"`
		Population       bool     `json:"population"`
		Area             bool     `json:"area"`
		TargetCurrencies []string `json:"targetCurrencies"`
	} `json:"features"`
	Lastchange string `json:"lastchange"`
}

type CurrentRegistration struct {
	Id   string `json:"id"`
	Time string `json:"time"`
}

type CountryRequest []struct {
	Capital    []string  `json:"capital"`
	Latlng     []float32 `json:"latlng"`
	Area       float32   `json:"area"`
	Population int       `json:"population"`
	Currencies map[string]struct {
	} `json:"currencies"`
}

type MeteoRequest struct {
	Daily struct {
		Temperature2MMax []float64 `json:"temperature_2m_max"`
		Temperature2MMin []int     `json:"temperature_2m_min"`
		PrecipitationSum []float64 `json:"precipitation_sum"`
	} `json:"daily"`
}

type DashboardResponse struct {
	Country  string `json:"country"`
	ISOcode  string `json:"isoCode"`
	Features struct {
		Temperature   float32 `json:"temperature"`
		Precipitation float32 `json:"precipitation"`
		Capital       string  `json:"capital"`
		Coordinates   struct {
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"langitude"`
		} `json:"coordinates"`
		Population       int                `json:"population"`
		Area             float32            `json:"area"`
		TargetCurrencies map[string]float32 `json:"targetcurrencies"`
	} `json:"features"`
	Lastchange string `json:"lastchange"`
}

type CurrencyRequest struct {
	Rates map[string]float32 `json:"rates"`
}

type Apendix []struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Cca2 string `json:"cca2"`
}

type Notification struct {
	Id      string `json:"id"`
	Url     string `json:"url"`
	Country string `json:"country"`
	Event   string `json:"event"`
}

type CurrentNotification struct {
	Id string `json:"id"`
}

type WebhookInvoke struct {
	Id      string `json:"id"`
	Country string `json:"country"`
	Event   string `json:"event"`
	Time    string `json:"time"`
}

type NotificationFirebase struct {
	Id      string `json:"id"`
	Country string `json:"country"`
	Event   string `json:"event"`
	Time    string `json:"time"`
	Url     string `json:"url"`
}

type Status struct {
	ContriesApi int     `json:"countries_api"`
	MeteoApi    int     `json:"meteo_api"`
	CurrencyApi int     `json:"currency_api"`
	Webhooks    int     `json:"webhooks"`
	Version     string  `json:"version"`
	UpTime      float64 `json:"uptime"`
}
