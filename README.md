# discord-bot

I wanted to experiment controlling my "smart" lights at home with a chat bot interface.

Originally wanted to use Google Chat for this task (as I've created a bot previously). Looks like Google Chat using Gmail does not support bots yet!

I'm very impressed with the API that Discord has created for bots, it's so well thought out.

Followed this tutorial / inspiration from <https://dev.to/aurelievache/learning-go-by-examples-part-4-create-a-bot-for-discord-in-go-43cf>.

# Things to tidy up

* Clean up error handling.
* Clean up signal handling.
* Figure out a clever way of allowing people other than adminUsername to interact with bot.
* Figure out how to abstract the hardcoded lights() func, and support other devices.