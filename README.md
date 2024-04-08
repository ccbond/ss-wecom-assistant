# SS Wecom Assistant

## WeChat Work Smart Customer Service with OpenAI Assistant

This project integrates the OpenAI Assistant API with WeChat Work (企业微信) to provide an intelligent customer service bot. Built with Go 1.21, it leverages the robustness of the [go-openai](https://github.com/sashabaranov/go-openai) library and the flexibility of [PowerWeChat](https://github.com/ArtisanCloud/PowerWeChat) to deliver a seamless and efficient customer service experience in corporate WeChat environments.

## Features

- **Intelligent Customer Interaction**: Utilizes OpenAI's cutting-edge language models for understanding and responding to customer queries.
- **Seamless Integration with WeChat Work**: Offers a fully integrated experience for WeChat Work users, providing smart responses directly in the app.
- **Customizable Response Handling**: Easily tailor the bot's responses to fit your organization's tone and guidelines.
- **Real-Time Analytics**: Monitor interactions and gather insights to continually improve the customer service experience.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- A valid OpenAI API key
- Access to WeChat Work with administrative privileges

### Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/ccbond/ss-wecom-assistant.git
   ```

2. Navigate to the project directory:

   ```bash
   cd ss-wecom-assistant
   ```

3. Build directly

   ```bash
   go build
   ```

4. Configure your OpenAI and WeChat Work API keys in `.env`:

   ```
   WECHAT_APP_ID=
   WECHAT_APP_SECRET=
   WECHAT_TOKEN=
   WECHAT_ENCODEING_AES_KEY=
   WECHAT_AGENT_ID=1
   WECHAT_ZJKFID=
   OPENAI_API_KEY=
   OPENAI_ASSISTANT_ID=
   ADMIN_SECRET=
   ```

5. Running the Server

   ```
   go build && ./ss-wecom-assistant or go run main.go
   ```

6. Usage
   To interact with the bot, send a message to the WeChat Work account linked with your application. The bot will process the message using the OpenAI API and respond accordingly.

7. FAQ
   Q: How do I obtain an OpenAI API key?
   A: You can apply for an API key by creating an account on the OpenAI website and following their API access guidelines.

Q: Can I customize the bot's responses?
A: Yes, response behavior can be customized by editing the response templates in the configuration.
