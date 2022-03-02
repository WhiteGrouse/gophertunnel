package main

import (
  "github.com/sandertv/gophertunnel/minecraft"
  "github.com/sandertv/gophertunnel/minecraft/auth"
  "github.com/sandertv/gophertunnel/minecraft/resource"
  "golang.org/x/oauth2"
  "os"
  "fmt"
  "io"
  "time"
)

func main() {
  token, err := auth.RequestLiveToken()
  if err != nil {
    panic(err)
  }
  src := auth.RefreshTokenSource(token)
  //IPv6で繋がるのを防ぐため
  servers := [] string {
    "mco.cubecraft.net:19132",
    "139.99.135.114:19132",//"mco.lbsg.net:19132",
    "51.79.234.180:19132",//"geo.hivebedrock.network:19132",
    "13.95.127.174:19132",//"play.pixelparadise.gg:19132",
    "51.89.217.217:19132",//"play.galaxite.net:19132",
    "52.166.239.43:19132",//"play.inpvp.net:19132",
    "108.178.12.50:19132",//"mco.mineplex.com:19132",
  }
  for _, server := range servers {
    download(server, src)
    println()
  }
}

func download(remote string, src oauth2.TokenSource) {
  fmt.Printf("Connecting to %s...\n", remote)
  seconds,  _ := time.ParseDuration("59s")
  serverConn, err := minecraft.Dialer {
    TokenSource: src,
  }.DialTimeout("raknet", remote, seconds)
  if err != nil {
    println("Error when Dial")
    fmt.Printf("%v\n", err)
    return
    //panic(err)
  }
  if err := serverConn.DoSpawn(); err != nil {
    println("Error when DoSpawn")
    fmt.Printf("%v\n", err)
    return
    //panic(err)
  }
  for _, pack := range serverConn.ResourcePacks() {
    save_pack(pack)
    fmt.Printf("Saved %s\n", pack.UUID())
  }
  serverConn.Close()
}

func save_pack(pack *resource.Pack) {
  key_path := fmt.Sprintf("%s.txt", pack.UUID())
  f, err := os.Create(key_path)
  if err != nil {
    panic(err)
  }
  f.WriteString(pack.ContentKey())
  f.Close()

  pack_path := fmt.Sprintf("%s.zip", pack.UUID())
  f, err = os.Create(pack_path)
  if err != nil {
    panic(err)
  }
  _, err = io.Copy(f, pack.Content())
  if err != nil {
    panic(err)
  }
  f.Close()
}
