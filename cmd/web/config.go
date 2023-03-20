package main

type configuration struct {
	BCHAddress       string   `json:"BCHAddress"`
	MinimumDonation  float64  `json:"MinimumDonation"`
	MaxMessageChars  int      `json:"MaxMessageChars"`
	MaxNameChars     int      `json:"MaxNameChars"`
	WebViewUsername  string   `json:"WebViewUsername"`
	WebViewPassword  string   `json:"WebViewPassword"`
	OBSWidgetRefresh string   `json:"OBSWidgetRefresh"`
	Checked          bool     `json:"ShowAmountCheckedByDefault"`
	EnableEmail      bool     `json:"EnableEmail"`
	SMTPServer       string   `json:"SMTPServer"`
	SMTPPort         string   `json:"SMTPPort"`
	SMTPUser         string   `json:"SMTPUser"`
	SMTPPass         string   `json:"SMTPPass"`
	SendToEmail      []string `json:"SendToEmail"`
}

func getDefaultConfig() configuration {
	var config configuration
	config.BCHAddress = "bitcoincash:address"
	config.MinimumDonation = 0.001
	config.MaxMessageChars = 300
	config.MaxNameChars = 25
	config.WebViewUsername = "admin"
	config.WebViewPassword = "adminadmin"
	config.OBSWidgetRefresh = "10"
	config.Checked = true
	config.EnableEmail = false
	config.SMTPServer = ""
	config.SMTPPort = ""
	config.SMTPUser = ""
	config.SMTPPass = ""

	return config
}
