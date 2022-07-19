# Deutschbot

Telegram bot for those who tries to learn German. Currently, it comes with [top 500 frequent German words](https://www.thegermanprofessor.com/top-500-german-words/).  

## Requirements
- sqlite (the data image is created at `bin/sqlite.db` so keep that in mind if you run the app in docker env)
- telegram token. You have to discuss it with [Botfather](https://t.me/botfather).  


## Launch

execution files are inside `cmd` directory.  
Initially you need to generate tasks to work with. You can do it by running  
```bash
go run cmd/tasks/generate/main.go
```

It may and will take a while due to communication with external translator API. 

After the tasks are successfully generated run the bot itself.  Don't forget to set environment variable  
`TG_BOT_TOKEN` (hope you got it from Botfather already).  

```bash
go run cmd/bot/main.go
```

## Commands

Some commands are stateful and will block others if the flow is in progress. For example [/start_quiz](#/start_quiz).  

- [/learn_word](#learn_word)
- [/start_quiz](#start_quiz)
    - [/example](#example)
- [/translate](#translate)

### /learn_word

Shows a random word

```
/learn_word
---

Question: How do you say "my"?
Applicable answers: mein.
Usages:
1. As a result, my daily routine became much more bearable.
Infolgedessen wurde mein Tagesablauf wesentlich erträglicher.
2. Now, we are at the end of my exemplary scenario.
Hiermit ist mein beispielhaftes Szenario abgeschlossen.
3. As a landscape and nature photographer, my tripod is my constant companion.
Als Landschafts- und Naturfotografin ist mein Stativ mein ständiger Begleiter.
```


### /start_quiz

Starts a quiz of ten questions to be solved.  Each comes separately.  

```
/start_quiz
---

QuizHandler started.
How do you say "child"?
---

ein Kind
---

Incorrect.Correct was: das Kind
How do you say "view"?
---

der Blick
---

Correct.
What does "bisschen" mean?
---

```

#### /example

This command is stateful and requires active quiz workflow. It will give some hints with a word usage.

```
Correct.
What does "bisschen" mean?
---

/example
---

Example
Ein bisschen Revolution wäre jetzt gut.

```

#### Stats

Every time you answer the question there is a statistic gathered. Next time when you start quiz it will attempt to use  
tasks that have worst total_attempts/correct_answers ratio. It won't take those tasks if there are not enough data gathered.  
Also, it won't take tasks that have you have given right answers last time they had appeared.

### /translate

Obviously translates some word from German to English. Experimental and non-reliable due to natural limitations of the  
translator that is used([reverso](https://context.reverso.net/)). For the same reason examples might be odd or even wrong.  

> **Please, do not abuse this feature. Respect each other.**

```
/translate Vater
---

Translation: father; dad; daddy
Usages:
1. Your father claims, quite strenuously, that he's innocent.
Ihr Vater beteuert nachdrücklich seine Unschuld.
2. Frieda's father had taken early retirement because of his war wound.
Friedas Vater war auf Grund seiner Kriegsverletzung Frührentner.
3. Stephan Noller is a certified psychologist, digital entrepreneur and dad of four daughters.
Stephan Noller ist Diplom-Psychologe, Digital-Unternehmer und Vater von vier Töchtern.
4. A toddler holding on to her dad's leg.
Ein Kleinkind Festhalten am Bein ihres Vater.
5. My father used to work in a building company.
Mein Vater arbeitete in einem Bauunternehmen.
```
