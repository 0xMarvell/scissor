# Scissor ✂️

Welcome to Scissor, the cutting-edge URL shortening API service! With Scissor, you can easily transform long, complex URLs into concise, user-friendly links. Whether you're building a web application, managing social media campaigns, or simply looking to share links efficiently, this API has got you covered. Say goodbye to unwieldy URLs and hello to streamlined communication. Let Scissor help you trim the excess and make every character count. Get started today and simplify your URL management with Scissor!

## Docs

The API documentation can be found [here](https://sci-ssor.onrender.com/docs/index.html)

## Run Scissor Locally

- Install [Docker and Docker Compose](https://docs.docker.com/engine/install/)
- In `pkg/config/database.go`, uncomment line `19` and comment out line `20`
- In `pkg/controllers/shortenControllers.go`, uncomment lines `63` and `67` and comment out lines `64` and `68`.
- Build Docker service

    ```bash
    docker compose up
    ```
