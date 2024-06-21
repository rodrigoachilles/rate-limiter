Traductions:

* [Anglais](README.md)
* [Portugais (Brésil)](README_pt_br.md)

---

# 📶 Limiteur de Taux (rate-limiter)

![Logo du Projet](assets/rate_limiter-logo.png)

Bienvenue dans le système de limiteur de taux développé en Go ! Ce projet vous permet de limiter le nombre de requêtes par seconde en fonction d'une adresse IP spécifique ou d'un jeton d'accès.

## 📑&nbsp;Table des Matières

- [📖 Introduction](#introduction)
- [🛠 Prérequis](#prérequis)
- [⚙️ Installation](#installation)
- [🚀 Utilisation](#utilisation)
- [🔍 Exemples](#exemples)
- [🤝 Contribution](#contribution)
- [📜 Licence](#licence)

## 📖&nbsp;Introduction

Ce système de limiteur de taux est un projet développé en Go qui permet de limiter le nombre de requêtes par seconde en fonction d'une adresse IP ou d'un jeton d'accès. Il aide à contrôler efficacement le trafic vers un service web.

## 🛠&nbsp;Prérequis

Assurez-vous d'avoir les éléments suivants installés avant de continuer :

- [Go](https://golang.org/doc/install)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## ⚙️&nbsp;Installation

1. Clonez ce dépôt :

    ```sh
    git clone git@github.com:rodrigoachilles/rate-limiter.git
    cd rate-limiter
    ```

2. Exécutez Docker Compose :

    ```sh
    docker-compose up -d
    ```

## 🚀&nbsp;Utilisation

Après avoir démarré Docker Compose, vous pouvez configurer et utiliser le limiteur de taux.

### 🔧&nbsp;Configuration

1. Créez un fichier `.env` à la racine du projet avec les configurations suivantes :

    ```env
   SERVER_PORT=:8080
   REDIS_ADDR=redis:6379
   IP_LIMIT=10
   TOKEN_LIMIT=100
   BLOCK_TIME=300 # secondes
    ```

2. Exécutez le serveur Go :

    ```sh
    go run main.go
    ```

### 📚&nbsp;Middleware

Pour utiliser le middleware du limiteur de taux, ajoutez-le à votre serveur HTTP :

```go
package main

import (
   "github.com/rodrigoachilles/rate-limiter/configs"
   "github.com/rodrigoachilles/rate-limiter/internal/infra/middleware"
   "github.com/rodrigoachilles/rate-limiter/internal/usecase/limiter"
   "github.com/rodrigoachilles/rate-limiter/internal/infra/database"
   "log"
   "net/http"
   "time"
   "context"
)

func main() {
   cfg, err := configs.LoadConfig()
   if err != nil {
      log.Fatalf("failed to load config: %v", err)
   }

   repo := repository.NewRedisRepository(cfg.RedisAddr)
   l := limiter.NewLimiter(repo, int64(cfg.IPLimit), int64(cfg.TokenLimit), cfg.BlockTime)

   mux := http.NewServeMux()
   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      _, _ = w.Write([]byte("Hello, World!"))
   })

   handler := middleware.RateLimiter(l)(mux)

   srv := &http.Server{
      Addr:         cfg.ServerPort,
      BaseContext:  func(_ net.Listener) context.Context { return ctx },
      WriteTimeout: 15 * time.Second,
      ReadTimeout:  15 * time.Second,
      Handler:      handler,
   }

   log.Println("Server is running on port 8080")
   log.Fatal(srv.ListenAndServe())
}
```

## 🔍&nbsp;Exemples

Voici quelques exemples d'utilisation du limiteur de taux :

- Limiter le nombre de requêtes par seconde en fonction de l'IP.
- Limiter le nombre de requêtes par seconde en fonction d'un jeton d'accès.
- Bloquer les nouvelles requêtes après avoir dépassé la limite pour un temps spécifié.

## 🤝&nbsp;Contribution

N'hésitez pas à ouvrir des issues ou à soumettre des "pull requests" pour des améliorations et des corrections de bugs.

## 📜&nbsp;Licence

Ce projet est sous licence MIT.
