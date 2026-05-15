# 🍹 PbxSpam Go

A fast, multi-bot Telegram spam and raid bot built in Go using the [gogram](https://github.com/amarnathcjd/gogram) library.

---

## ⚡ Features

- Multi-bot support (run multiple bots at once)
- Spam commands with reply support
- Multi-language raid commands (Hindi, English, Punjabi)
- Reply raids — tag a specific message repeatedly
- Shayari command — send random Punjabi shayaris
- MongoDB persistence — active raids resume after bot restart
- Sudo-only commands — only the owner can use powerful commands

---

## 🛠️ Requirements

- Go 1.21+
- MongoDB (optional but recommended for raid persistence)
- Telegram API credentials from [my.telegram.org](https://my.telegram.org)

---

## 🚀 Setup

**1. Clone the repo**
```bash
git clone https://github.com/badmunda05/SpamGo
cd SpamGo
```

**2. Create your `.env` file**
```bash
cp sample.env .env
nano .env
```

**3. Fill in your `.env`**
```env
APP_ID=12345678
APP_HASH=your_app_hash_here
OWNER_ID=your_telegram_user_id
MONGO_URL=mongodb://localhost:27017   # optional

BOT_TOKEN1=your_first_bot_token
BOT_TOKEN2=your_second_bot_token      # add more as needed

START_PIC=https://your-image-url.jpg  # optional
HELP_PIC=https://your-image-url.jpg   # optional
```

**4. Run**
```bash
go run main.go
```

---

## 📋 Commands

### General
| Command | Description |
|---------|-------------|
| `/start` | Show bot info |
| `/help` | Show help menu with buttons |
| `.ping` | Check bot response speed |
| `.restart` | Restart the bot (sudo only) |

---

### 💬 Spam
| Command | Example | Description |
|---------|---------|-------------|
| `.spam` | `.spam 10 hello` | Send a message N times |
| `.ds` | `.ds 1 10 hello` | Delay spam — send with custom delay (seconds) |
| `.delayspam` | `.delayspam 1 10 hello` | Same as `.ds` |
| `.sspam` | `.sspam 10` | Spam a sticker/media (reply to it) |

> **Tip:** If you use `.spam` or `.ds` while replying to a message, the bot will tag that message each time.

---

### 🔥 Raid
| Command | Language | Description |
|---------|----------|-------------|
| `.raid N` | Bold/Stylish | Random normal raid |
| `.hraid N` | Hindi | Random Hindi raid |
| `.eraid N` | English | Random English raid |
| `.punraid N` | Punjabi | Random Punjabi raid |
| `.replyraid N` | Bold/Stylish | Reply raid (reply to a msg first) |
| `.hreplyraid N` | Hindi | Hindi reply raid |
| `.ereplyraid N` | English | English reply raid |
| `.preplyraid N` | Punjabi | Punjabi reply raid |

> **Note:** Active raids are saved to MongoDB. If the bot restarts during a raid, it will automatically resume.

---

### 🌹 Shayari
| Command | Example | Description |
|---------|---------|-------------|
| `.shayari` | `.shayari 2` | Reply to a user's message to send N random Punjabi shayaris as replies |

---

## 📁 Project Structure

```
SpamGo/
├── main.go              # Entry point
├── .env                 # Your config (not committed)
├── sample.env           # Example config
├── client/
│   └── client.go        # Bot init, login, handler registration
├── config/
│   ├── config.go        # Env loading
│   └── data.go          # Raid lists, shayari list
├── database/
│   ├── db.go            # MongoDB connection
│   ├── sudo.go          # Sudo user management
│   └── raid.go          # Active raid persistence
└── modules/
    ├── module.go         # Command/callback registration system
    ├── utils.go          # Helper functions
    ├── basic.go          # /start, /help, .ping, .restart
    ├── spam.go           # .spam, .ds, .sspam
    ├── raid.go           # All raid commands
    └── shayari.go        # .shayari command
```

---

## 👤 Developer

- **Developer:** [@BadmundaXd](https://t.me/BadmundaXd)
- **Channel:** [@PBX_UPDATE](https://t.me/PBX_UPDATE)
- **Support:** [@PBXCHATS](https://t.me/PBXCHATS)
- **Repo:** [github.com/badmunda05](https://github.com/badmunda05)

---

## 📄 License

MIT License — feel free to use and modify.
