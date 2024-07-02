package main_test

import (
	"testing"

	"github.com/IceWhaleTech/CasaOS-AppManagement/cmd/validator/pkg"
	utils_logger "github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/stretchr/testify/assert"
)

func TestAppIsCorrect(t *testing.T) {
	utils_logger.LogInitConsoleOnly()

	err := pkg.VaildDockerCompose([]byte(CorrectYAML))
	assert.NoError(t, err)
}

func TestTipsIsWrong(t *testing.T) {
	utils_logger.LogInitConsoleOnly()

	err := pkg.VaildDockerCompose([]byte(TipsWrongYAML))
	assert.NotEmpty(t, err)
}

const (
	CorrectYAML   = ``
	TipsWrongYAML = `name: swingmusic
services:
  swingmusic:
    image: ghcr.io/swingmx/swingmusic:v1.4.8
    container_name: swingmusic
    volumes:
      - type: bind
        source: /DATA/AppData/$AppID/config
        target: /config
      - type: bind
        source: /DATA/Media/Music
        target: /music
    ports:
      - "1970:1970"
    restart: unless-stopped
x-casaos:
  architectures:
    - amd64
  main: swingmusic
  author: CasaOS Team
  category: Media
  description:
    en_us: |
      Swing Music is a beautifully designed, self-hosted music streaming server. Like a cooler Spotify ... but bring your own music.
    zh_cn: |
      Swing Music 是一款设计精美的自托管音乐流媒体服务器。就像更酷的 Spotify......但可以带来你自己的音乐。
  developer: SwingMX
  icon: https://cdn.jsdelivr.net/gh/IceWhaleTech/CasaOS-AppStore@main/Apps/SwingMusic/icon.png
  screenshot_link:
  - https://cdn.jsdelivr.net/gh/IceWhaleTech/CasaOS-AppStore@main/Apps/SwingMusic/screenshot-1.png
  - https://cdn.jsdelivr.net/gh/IceWhaleTech/CasaOS-AppStore@main/Apps/SwingMusic/screenshot-2.png
  - https://cdn.jsdelivr.net/gh/IceWhaleTech/CasaOS-AppStore@main/Apps/SwingMusic/screenshot-3.png
  - https://cdn.jsdelivr.net/gh/IceWhaleTech/CasaOS-AppStore@main/Apps/SwingMusic/screenshot-4.png
  - https://cdn.jsdelivr.net/gh/IceWhaleTech/CasaOS-AppStore@main/Apps/SwingMusic/screenshot-5.png

  tagline:
    en_us: Swing Music is a beautifully designed, self-hosted music streaming server. Like a cooler Spotify ... but bring your own music.
    zh_cn: Swing Music 是一款设计精美的自托管音乐流媒体服务器。就像更酷的 Spotify......但可以带来你自己的音乐。
  tips: |
    before_install:
      en_us: |
        When you first start Swing Music, it will ask you to pick music directory: Where do you want to look for music?
        select "Specific directories" Option, and select "/music" and rescan.
      zh_cn: |
        首次启动 Swing Music 时，它会要求您选择音乐目录：您想在哪里查找音乐？
        选择“特定目录”选项，然后选择“/music”并重新扫描。
  title:
    en_us: Swing Music
    zh_cn: Swing Music
  index: /
  port_map: "1970"
`
)
