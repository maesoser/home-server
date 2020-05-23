---
title: Quantenverschränkung, la circustancia que provocará el efecto 2000 del futuro.
draft: false
tags: networks, engineering
---

¿Os acordáis del efecto 2000? Ese bug apocalíptico que según la tele estuvo a punto de poner en jaque nuestros sistemas informáticos la noche del 31 de diciembre de 1999.

[![Efecto2000](http://img.youtube.com/vi/Lf2eT8kmunE/0.jpg)](http://www.youtube.com/watch?v=Lf2eT8kmunE "")

El problema se daba en programas que se habían construido bajo una premisa específica, que dichos sistemas jamás iban a utilizarse más allá del siglo XX. Eso permitía a los programadores guardar el año en dos dígitos en vez de en cuatro, y obtener un ahorro de memoria considerable en aquella época.

Yo era bastante pequeño, pero recuerdo que al ver las noticias en la tele me dio la sensación de que nos habíamos dado cuenta ese mismo año de esa circunstancia tan peligrosa. Y lo cierto es que ya desde finales de los 80, [algunos incluso antes](https://www.independent.co.uk/news/obituaries/bob-bemer-550018.html), la industria del software había empezado a tomar medidas de cara al nuevo milenio.

Tendremos otros pequeños sustos en el futuro. Uno de ellos ocurrirá en el [año 2038](https://en.wikipedia.org/wiki/Year_2038_problem) cuando el contador de 32 bit que se usa en algunos sistemas operativos para guardar la fecha se desborde. Pero otro bastante más serio ~y mucho más molón~ podría tener lugar en el momento en el que nuestras comunicaciones dejen de usar la física convencional y empecemos a desplegar sistemas de comunicación basados en física cuántica.

Seguro que os suena escuchar de cuando en cuando alguna noticia sobre un equipo de física que ha [logrado teletransportar algo](https://www.technologyreview.com/s/608252/first-object-teleported-from-earth-to-orbit/). Detrás de esos titulares tan impactantes en los que yo siempre espero encontrarme una apasionante historia sobre teletransportar un conejo o un coche suelen esconderse artículos algo vagos y aburridos que en mi opinión no logran captar lo verdaderamente interesante del asunto. 

![star_trek_teleportation](https://media.giphy.com/media/nYaRWwyG9qAH6/giphy.gif)

Porque al final, tras leerlo te enteras que lo único que han conseguido es transmitir de de una partícula elemental a otra el momento de espín, o alguna otra propiedad física, extraña y desconocida para la mayoría de nosotros. ¿Qué tiene eso de utilidad? Y lo que nos interesa en este artículo ¿En qué puede afectar eso a nuestras telecomunicaciones?

Con lo que estos equipos de físicos están realmente experimentando es con una propiedad hipotética que tendrían las partículas elementales llamada entrelazamiento. Según esa propiedad, todas las partículas pertenecientes a un mismo sistema cuántico compartirían ciertas características físicas. Un cambio de una de esas características en una partícula de ese sistema provocará que todas las demás partículas cambien también. Y de forma **instantánea**.

Muy pocos protocolos en Internet están preparados para esta situación, en especial los [algoritmos de control de congestión](https://en.wikipedia.org/wiki/TCP_congestion_control#Compound_TCP) que se encargan de regular la velocidad a la que se transmite la información desde un ordenador para impedir saturar al receptor.

Algunos de ellos, como [LedBat](https://tools.ietf.org/html/rfc6817), intentan molestar lo menos posible al resto de comunicaciones de la red, teniendo en cuenta en sus cálculos el tiempo que pasa entre enviar un mensaje y recibir confirmación de que ha sido recibido. Esa medida se llama Round Trip Time (RTT). En una red cuántica, este RTT podría llegar a ser 0, lo que provocaría que LedBat enviara más tráfico del que la red pudiera soportar. LedBat se usa en la actualidad en la red torrent y para descargar actualizaciones en segundo plano en algunos sistemas operativos.

Otro protocolo que vería aumentada su agresividad sería CompoundTCP, el algoritmo de control de congestión de Windows, o [BBR](https://blog.apnic.net/2020/01/10/when-to-use-and-not-use-bbr/), el prometedor algoritmo que usa Google para transmitir los videos de youtube.

Pero tranquilos, a pesar de que la cuenta atrás ya ha empezado, se cree que la gran mayoría del tráfico basado en TCP resistiría este cambio y seguiría funcionando, aunque no de forma óptima.

## Referencias

[The Quantum Bug](https://www.rfc-editor.org/rfc/rfc8774.txt)

[Design Considerations for Faster-Than-Light (FTL) Communication](https://tools.ietf.org/html/rfc6921)

https://www.youtube.com/watch?v=Vwvix4ahHEk

