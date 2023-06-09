# go-worker-command
Worker que lê mensagens de filas diferentes e processa fazendo requisições externas em paralelo

[![LinkedIn][linkedin-shield]][linkedin-url]


<!-- SOBRE O PROJETO -->
## Sobre o Projeto

Aplicação que lê mensagens em parapelo de duas filas com DQL e routingkeys diferentes utilizando canais e goroutines. Faz o processamento das mensagens chamando APIs externas em paralelo utilizando WaitGroup e goroutines 

![img](https://user-images.githubusercontent.com/21323326/233877399-487d793c-76b4-445b-88fd-111c94145c26.png)


### Utilizando

* [![Go][Go-badge]][Go-url]
* [![RabbitMQ][RabbitMQ-badge]][RabbitMQ-url]
* [![Docker][Docker-badge]][Docker-url]
* [![VS Code][VSCode-badge]][VSCode-url]


<!-- GETTING STARTED -->
## Getting Started

Instruções para execução da aplicação

### Prerequisites

Executar o comando para inicializar o MongoDB, RabbitMQ e a aplicação na porta selecionada
* docker
  ```sh
  docker-compose up -d
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/your_username_/Project-Name.git
   ```
2. exec
   ```sh
   go run main.go
   ```

<!-- AGRADECIMENTOS -->
## Agradecimentos

* [Finnhub](https://finnhub.io/)
* [AwesomeApis](https://docs.awesomeapi.com.br/api-de-moedas)

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/felipe-fernandes-fca/
[Go-url]: https://golang.org/
[Go-badge]: https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white
[MongoDB-badge]: https://img.shields.io/badge/mongodb-%234ea94b.svg?style=flat&logo=mongodb&logoColor=white
[MongoDB-url]: https://www.mongodb.com/
[RabbitMQ-badge]: https://img.shields.io/badge/rabbitmq-%23ff6600.svg?style=flat&logo=rabbitmq&logoColor=white
[RabbitMQ-url]: https://www.rabbitmq.com/
[Docker-badge]: https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[VSCode-badge]: https://img.shields.io/badge/VS_Code-007ACC?style=flat&logo=visual-studio-code&logoColor=white
[VSCode-url]: https://code.visualstudio.com/
