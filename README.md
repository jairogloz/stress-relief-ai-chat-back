# 🌟 Stress Relief AI Chat Backend

## 🧠 Purpose

The **Stress Relief AI Chat** is a specialized backend service designed to support a chat application with a single, powerful goal: **to make you feel better no matter what**. Unlike a general-purpose AI chat, this app is laser-focused on providing comfort, actionable advice, and a sense of calm in every interaction.

### 🎯 What Makes It Different?

- **Always Positive:** Every response is crafted to uplift and reassure.
- **Action-Oriented:** Not just words—every chat session includes **specific, actionable steps** you can take to improve your situation or your mood.
- **Real-Time AI Responses:** The app currently connects with **OpenAI's GPT-4**, with plans to integrate **additional AI models** in the future to ensure diverse and tailored responses.

---

## 🛠️ How It Works

The backend serves as a bridge between the frontend chat interface and various AI assistants. It manages user authentication through **Supabase** (currently under development) and handles AI interactions using the **github.com/sashabaranov/go-openai** library. The application follows a **Hexagonal Architecture** to ensure clean code and maintainability.

### 🔍 Key Features

1. **User Authentication (Under Development):** The app will use **Supabase** for easy signup and login.
2. **AI Chat Endpoint:** Seamless integration with the frontend to process user inputs and generate AI responses.
3. **Real-Time Support:** Currently connects directly to **OpenAI GPT-4** and aims to support additional AI models.
4. **No Chat Storage:** Ensures privacy by acting only as a bridge—no personal data or chat history is stored.

---

## 📬 Join the Early Access List!

Are you interested in trying out the **Stress Relief AI Chat**? Sign up for **early access** and be among the first to experience this unique tool. You'll also receive updates about new features and improvements.

👉 **[Sign Up for Early Access](https://ai-stress-relief-landing-page.vercel.app/)**

---

## 🚀 Getting Started (For Developers)

### Prerequisites

- **Golang** (v1.18 or higher)
- **Supabase Account** for authentication
- **OpenAI API Key**

### Setup Instructions

1. **Clone the Repository:**

```bash
git clone https://github.com/yourusername/stress-relief-ai-chat-backend.git
cd stress-relief-ai-chat-backend
```

2. **Install Dependencies:**

```bash
go mod tidy
```

3. **Set Up Environment Variables:**

```bash
cp .env.example .env
# Add your Supabase and OpenAI credentials to the .env file
```

4. **Run the Application:**

```bash
go run cmd/main.go
```

5. **Test the API:**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"message":"I feel stressed."}' http://localhost:8080/chat
```

---

## 🧑‍💻 Contributing

Contributions are welcome! Please open an **issue** or submit a **pull request** if you want to improve the project.

---

## 📄 License

This project is licensed under the **MIT License**.
