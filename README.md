## SMS Youtube

In a few days we are going to have a [meetup](http://www.meetup.com/IPAs-APIs/events/228942964/) at [the office](http://messagebird.com) and I thought of doing something silly so that the people attending could play videos and listen to music together (youtube party yeeyyy!!). 

So I made this little hack real quick where people could SMS a video title to a VMN(Virtual Mobile Number) and it would show on a page openned on a computer with a screen on the kitchen (that's where we do parties at MessageBird). 

## Wanna try?

 * goto https://smsyoutube.herokuapp.com/
 * sms "Never Gonna Give You Up" to +447860039713 (or +31635777839 :nl:) 
 * RickRoll'D!

## How does it work?

I got a [MessageBird VMN](https://www.messagebird.com/en/virtual-mobile-number) and configured it to do a **POST** to a url on every message.

In the end of that url there is a small **go** code (hosted on **heroku**) that does a search on youtube and inserts the video result id into a **firebase** list.

On the other end, on that cool html page, there is a javascript code that watches that firebase list and when there is a new entry it tells the youtube embed, via the **youtube iframe api**, to bluntly switch to that video. 


