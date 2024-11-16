# Random chat
Was created for communication. The chat is designed for 1 on 1 communication. It works on a room system.
The project is made using gorilla/websocket

# Setup
1. Rename .env.example to .env.
2. Update .env with your configuration values.

# **API Endpoints:**
1. # Home
    Endpoint: `/random/home`
    Method: GET

2. # Chat
    Endpoint: `/random/chat`
    Method: GET
    Description: The main chat page, where communication with other users is provided. The websocket connection is established automatically.

3. # Swagger Documentation
    Endpoint: `/doc/*any`
    Method: GET
    Description: View Swagger documentation for all available endpoints.

4. # WebSocket
    Endpoint: `/random/ws`
    Method: GET
    Description: The URL is provided for establishing a websocket connection. Also responsible for processing everything that happens in the chat based on a websocket connection.


# **Running the Service**
1. Start the database

    ```bash
    make start
    ```


# Please don’t look at the front-end part) GPT chat helped me write it and I’m not good at it. Sry(