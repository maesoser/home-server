---
title: Para qué hacerlo fácil si puedes complicarte sin motivo: El nacimiento de Blog
draft: true
tags: programming
---


Llevo varios años queriendo hacer un blog. Entre la vaguería y la desidia todos esos proyectos han acabado languideciendo conteniendo apenas un par de artículos mediocremente escritos. En este tercer intento me he propuesto empezar de forma diferente. Dándole vueltas sobre la forma ideal de comenzarlo, alguna manera para evitar no cumplir mi objetivo. Soy muy poco disciplinado y tampoco soy demasiado bueno redactando, pero no se me da mal programar. ¿Por qué no, antes de escribir un blog, *escribo* un blog?

Es por eso que, a pesar de la [multitud de opciones disponibles](https://www.staticgen.com) para generar y servir una página web estática, he decidido crear mi pequeña utilidad para esta tarea: Blog.

#### Blog es simple 

Blog tan sólo necesita un fichero de configuración (`config.yaml`) y una carpeta con posts. El fichero de configuración incluye algunos detalles para personalizar el blog y la localización de la carpeta de posts:

```
---
title: Tidy Bits
subtitle: The blog where I tell my geek adventures
author: Maesoser
author_url: https://github.com/maesoser
url: blog.souvlaki.cf
year: '2020'
posts: "./posts"
```

En la esa carpeta deberían encontrarse los distintos artículos, ordenados como subcarpetas.

```
  posts
    ├── 010120_post1
    │   ├── foto1.jpg
    │   └── post.md
    ├── 120220_quantumbug
        ├── 120220_quantumbug.md
        ├── foto2.jpg
        └── foto4.png
```

Iniciar Blog es tan sencillo como:

```
blog -port 8080 -tls_port 8443 -config /temp/config.yaml
```

![blog container starting](blog_start.png)

#### Blog es limpio

Soy un desastre con el diseño. Para conseguir que el blog fuera lo más sencillo posible a la par que agradable y simple de leer, he hecho uso de [Marx, un CSS sin clases](https://mblode.github.io/marx/). Esto me permite aplicarlo a mi código HTML sin tener que modificar cada elemento sobre el que quiera aplicar un formato. Añadirlo (y personalizarlo) es tan sencillo como incluir estas líneas en la cabecera de tu web:

```
<link rel='stylesheet' href='https://rawgit.com/mblode/marx/master/css/marx.min.css'>
<style> 
    main {  max-width: 1200px; } 
    img { display: block; margin-left: auto; margin-right: auto; }
</style>
```

Para colorear los fragmentos de código que añado a mis artículos uso [Highlight.js](https://highlightjs.org/). Añadirlo es igual de sencillo que con Marx. Tan sólo hace falta incluir unas pocas lineas a la cabecera de nuestra plantilla:

```
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.0.3/styles/tomorrow.min.css">
<style> .hljs { background: #F6F8FA; } </style>
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.0.3/highlight.min.js"></script>
<script>hljs.initHighlightingOnLoad();</script>
```

#### Blog es Seguro

Blog escucha por defecto tanto en el puerto 80 como en el puerto 443. Si no se ha incluido un certificado SSL, Blog emitirá un certificado firmado por él mismo al iniciarse.

Blog también sigue [las recomendaciones básicas](https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/) a la hora de exponer un servicio escrito en Go en Internet

## Desplegando la solución

Para desplegar el binario en mi servidor he hecho uso de Docker. Con un Dockerfile sencillo como el siguiente, la imagen final no llega a 10Mb de tamaño. Eso se debe a que hemos usado la imagen [scratch](https://docs.docker.com/develop/develop-images/baseimages/#create-a-simple-parent-image-using-scratch) como base, que tan sólo incluye el binario que hemos generado. Este truco no sólo nos ayuda a reducir espacio si no que también mejora la seguridad al deshacernos de todas los bnarios y objetos no necesarios dentro del contenedor. 

```
FROM golang:alpine as builder

RUN adduser -D blog
WORKDIR /
COPY . .

RUN apk add git && \
  go get -u github.com/gorilla/mux && \
  go get github.com/k3a/html2text && \
  go get github.com/gomarkdown/markdown && \
  go get -u github.com/go-bindata/go-bindata/... && \
  go-bindata -pkg main -o assets.go assets/ && \
  CGO_ENABLED=0 GOOS=linux go build -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o blog .

FROM scratch
COPY --from=builder /blog /blog
COPY --from=0 /etc/passwd /etc/passwd
USER blog
WORKDIR /
ENTRYPOINT ["./blog"]
```

Un fichero de docker-compose igual de simple configura el contenedor:

```
version: "2.4"

services:
  blog:
    container_name: blog
    restart: unless-stopped
    mem_limit: 64m
    build:
     context: ./containers/blog
     dockerfile: Dockerfile
    environment:
     - SERVER_HTTP_PORT=8080
     - SERVER_HTTPS_PORT=8443
    volumes:
     - ./config/blog/config.json:/config.yml
     - ./config/blog/posts:/posts
```

## Referencias

[Begginers guide to serving files using go](https://medium.com/rungo/beginners-guide-to-serving-files-using-http-servers-in-go-4e542e628eac)

[HTTP2 and TLS client and server example](http://www.inanzzz.com/index.php/post/9ats/http2-and-tls-client-and-server-example-with-golang)
