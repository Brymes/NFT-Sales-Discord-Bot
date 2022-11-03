<div align="center">
    <hr>
    <h1>NFT Sales Tracker</h1>
    <strong>
        A simple discord bot for tracking sales from specific NFT collections as well as tracking sales via Thresholds and other utility information
    </strong><br><br>
    <a href="https://bit.ly/3rDTqSW"><img src="https://img.shields.io/badge/%20-INVITE%20BOT-7F00FF.svg?style=for-the-badge&logo=discord" height="45" /></a>
<br>
</div>



---

# Intro

A simple Discord bot for fetching useful information about NFTs powered by [**DIA**](https://www.diadata.org).

### Github

https://github.com/Brymes/NFT-Sales-Discord-Bot

---

# Features

- Subscribe to Sales Updates for any NFT Collection
- Subscribe to Sales updates for all NFT Sales above a specified price
- Retrieve Floor price and useful price insights(e.g Moving Average, 24h Volume) of any NFT collection
- Personalize by Deploying on your own server with Heroku(WIP) or Docker instructions down below

Slash Commands
---

| Commands        | Description                                                                                                                                                            |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| /help           | Returns All Commands and their corresponding Descriptions                                                                                                              |
| /subscriptions  | Returns a list of commands    which the server has enabled                                                                                                             |
| /sales          | Allow to set-up a channel that receives all NFT collection sales matching the supplied contract address from the DIA NFT Event WebSocket to a selected discord channel | 
| /sales_stop     | Stops Bot from pushing sales update from a contract address or stop all Contract Address tracking bots if no contract address was provided                             |                                                                                                                                                                  |
| /floor          | Return floor price of the provided NFT collection contract address                                                                                                     |
| /all_sales      | Return all sales above the provided threshold to selected channel                                                                                                      |
| /all_sales_stop | Stop bot for all sales above the predetermined threshold and contract address                                                                                          |
| /stop_all       | Stops all bots from operating in the selected channel or stop all bots if channel is not provided                                                                      |

---

# Invite

Click this button the bot to your Server

<a href="https://bit.ly/3rDTqSW"><img src="https://img.shields.io/badge/%20-INVITE%20BOT-7F00FF.svg?style=for-the-badge&logo=discord" height="35" /></a>

# Docker

This section is for devs or If you intend on hosting this app yourself for better performance

### Prerequisites

- A Server with Docker installed e.g. [**Ubuntu Example**](https://docs.docker.com/engine/install/ubuntu/).
- A Discord Bot Token as [**described
  here**](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token).
- A Postgres Database as described [here](https://www.makeuseof.com/install-configure-postgresql-on-ubuntu/). 
- A channel to receive any uncaught errors. Read [**here**](https://turbofuture.com/internet/Discord-Channel-ID) on how
  to get a channel ID.

### Command to Run

Run the following command. Do replace the required variables with your sourced pre-requisites without the tags

```
sudo docker run -d -t -i --name my-nft-sales-tracker -it --net=host \
-e DISCORD_BOT_TOKEN="<YOUR_DISCORD_BOT_TOKEN>" \
-e DB_HOST="<YOUR_DB_HOST>" \
-e DB_USERNAME="<YOUR_DB_USERNAME>" \
-e DB_PASSWORD="<YOUR_DB_PASSWORD>" \
-e DB_DATABASE="<YOUR_DB_DATABASE>" \
-e DB_PORT="<YOUR_DB_PORT>" \
-e PANIC_CHANNEL="<YOUR_PANIC_CHANNEL_ID>" \
brymes/discord-nft-sales-tracker:stable 
```