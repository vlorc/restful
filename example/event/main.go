package main

import (
	"context"
	"github.com/vlorc/restful/pkg/engine"
	_ "github.com/vlorc/restful/pkg/engine/chi"
	"github.com/vlorc/restful/pkg/render"
	"github.com/vlorc/restful/pkg/web"
	"log"
	"net/http"
	"time"
)

var page = web.HTML(`
<!DOCTYPE html><html><body><h1>Getting server updates</h1><div id="result"></div><script>
if(typeof(EventSource) !== "undefined") {
  var source = new EventSource("/message");
  source.onopen = function(){
	console.log("opening");
  }
  source.onclose = function(){
	console.log("closed");
  }
  source.onerror = function(e){
	console.log("error:",e);
  }
  source.onmessage = function(e) {
	console.log(e)
    document.getElementById("result").innerHTML += e.data + "<br>";
  };
} else {
  document.getElementById("result").innerHTML = "Sorry, your browser does not support server-sent events...";
}
</script></body></html>
`)

func main() {
	g := engine.Default()

	q := make(chan string)

	go func() {
		for range time.NewTicker(time.Second * 2).C {
			select {
			case q <- time.Now().Format("2006-01-02 15:04:05"):
				log.Println("send ok~")
			default:
			}
		}
	}()

	g.Get("/message", func() func(context.Context) render.Render {
		return func(ctx context.Context) render.Render {
			select {
			case <-ctx.Done():
			case msg, ok := <-q:
				if ok {
					return render.Event{
						Data: msg,
					}
				}
			}
			log.Println("message exit")
			return nil
		}
	})

	g.Get("/page", page)

	http.ListenAndServe(":1234", g)
}
