package sailor

import (
	"github.com/astaxie/beego"
        "runtime"
)

func BeegoInit() {
        runtime.GOMAXPROCS(8)
        // beego.SetLogger("console", "")
	beego.SetLogger("file", `{"filename":"` + beego.AppConfig.String("logfile") + `", "maxlines":1000000}`)
        beego.SetLogFuncCall(true)
        beego.SetLevel(beego.LevelInformational)
        beego.Info( "appname:", beego.AppConfig.String("appname"), ", runmode:", beego.BConfig.RunMode )
}
