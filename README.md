# SlackMeme

Generate Memes using the ImgFlip API directly in Slack.

### Setup

Example of setting up with Heroku.

```
git clone git@github.com:adamgeorgeson/slack_meme.git
cd slack_meme
heroku create -b https://github.com/kr/heroku-buildpack-go.git
```

Manage your Slack Integrations.
 - Create a Slash Command with `/meme` setting the URL to the Heroku deployment.
 - Create an Incoming Webhook.

Register with ImgFlip.

```
heroku config:set SLACK_TOKEN=<token from Slash Command'
heroku config:set SLACK_WEBHOOK=<webhook url from Slash Incoming Webhook>
heroku config:set IMGFLIP_USERNAME=<imgflip_username>
heroku config:set IMGFLIP_PASSWORD=<imgflip_password>
git push heroku master
```
### Usage

```
/meme list # Returns the ImgFlip URL for top 100 Meme template IDs
/meme 405658 top: Another Slack Bot? bottom: I Hate It
```
