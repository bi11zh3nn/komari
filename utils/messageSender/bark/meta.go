package bark

type Addition struct {
	ServerURL string `json:"server_url" name:"Bark server URL" required:"true" default:"https://api.day.app" help:"Bark server URL, such as https://api.day.app"`
	DeviceKey string `json:"device_key" name:"Device key" required:"true" help:"Device key shown in the Bark app"`
	Icon      string `json:"icon" name:"Icon" help:"Push notification icon URL or system icon name"`
	Level     string `json:"level" name:"Level" type:"option" default:"timeSensitive" options:"active,timeSensitive,passive,critical"`
}
