import { useEffect, useState } from "react";
import "./App.css";

const API = import.meta.env.VITE_API_URL || "http://localhost:8080";

function App() {
  const [mode, setMode] = useState("login");
  const [token, setToken] = useState(localStorage.getItem("token") || "");

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);

  const [pollId, setPollId] = useState(
    new URLSearchParams(window.location.search).get("poll") || ""
  );

  const [poll, setPoll] = useState(null);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const isLoggedIn = !!token;

  async function handleAuth(e) {
    e.preventDefault();
    setMessage("");
    setLoading(true);

    const endpoint = mode === "login" ? "/auth/login" : "/auth/register";

    const body =
      mode === "login"
        ? { email, password }
        : { name, email, password };

    try {
      const response = await fetch(API + endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Something went wrong");
      }

      if (mode === "register") {
        setMessage("Registration successful. Please login.");
        setMode("login");
      } else {
        localStorage.setItem("token", data.token);
        setToken(data.token);
        setMessage("Login successful!");
      }
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function createPoll(e) {
    e.preventDefault();
    setMessage("");

    const cleanOptions = options.filter((item) => item.trim() !== "");

    if (cleanOptions.length < 2) {
      setMessage("Please provide at least 2 options.");
      return;
    }

    try {
      const response = await fetch(API + "/polls", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          question,
          options: cleanOptions,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Could not create poll");
      }

      const id = data.poll.id;
      setPollId(id);

      window.history.pushState({}, "", `/?poll=${id}`);

      setMessage("Poll created successfully!");
      setQuestion("");
      setOptions(["", ""]);

      loadPoll(id);
    } catch (error) {
      setMessage(error.message);
    }
  }

  async function loadPoll(id = pollId) {
    if (!id) return;

    try {
      const response = await fetch(`${API}/polls/${id}`);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Poll not found");
      }

      setPoll(data);
    } catch (error) {
      setMessage(error.message);
    }
  }

  async function vote(optionId) {
    if (!pollId) return;

    try {
      const response = await fetch(`${API}/polls/${pollId}/vote`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ optionId }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Vote failed");
      }

      setPoll((current) => {
        if (!current) return current;

        return {
          ...current,
          options: current.options.map((option) =>
            option.id === optionId
              ? { ...option, count: data.count }
              : option
          ),
        };
      });
    } catch (error) {
      setMessage(error.message);
    }
  }

  function logout() {
    localStorage.removeItem("token");
    setToken("");
    setMessage("Logged out.");
  }

  function addOption() {
    if (options.length < 6) {
      setOptions([...options, ""]);
    }
  }

  function updateOption(index, value) {
    const updated = [...options];
    updated[index] = value;
    setOptions(updated);
  }

  useEffect(() => {
    if (pollId) {
      loadPoll(pollId);
    }
  }, [pollId]);

  // WebSocket live updates
  useEffect(() => {
    if (!pollId) return;

    const WS_URL = API.replace(/^http/, "ws");

const ws = new WebSocket(
  `${WS_URL}/ws/polls/${pollId}`
);

    ws.onmessage = (event) => {
      try {
        const update = JSON.parse(event.data);

        setPoll((current) => {
          if (!current) return current;

          return {
            ...current,
            options: current.options.map((option) =>
              option.id === update.optionId
                ? { ...option, count: update.count }
                : option
            ),
          };
        });
      } catch {
        console.log("Invalid live update");
      }
    };

    ws.onerror = () => {
      console.log("WebSocket connection error");
    };

    return () => ws.close();
  }, [pollId]);

  const totalVotes =
    poll?.options?.reduce((sum, option) => sum + Number(option.count || 0), 0) ||
    0;

  return (
    <div className="app">
      <header className="navbar">
        <div className="logo">LivePoll</div>

        <div>
          {isLoggedIn ? (
            <button className="outline-btn" onClick={logout}>
              Logout
            </button>
          ) : (
            <button
              className="outline-btn"
              onClick={() =>
                setMode(mode === "login" ? "register" : "login")
              }
            >
              {mode === "login" ? "Sign up" : "Login"}
            </button>
          )}
        </div>
      </header>

      <main className="container">
        {message && <div className="message">{message}</div>}



        {!isLoggedIn && !pollId && (
  <section className="auth-layout">
    <div className="hero-panel">
      <div className="hero-content">
        <div className="hero-badge">
          <span className="hero-dot"></span>
          REAL-TIME POLLING
        </div>

        <h1 className="hero-title">
          Ask.
          <br />
          <span>Vote.</span>
          <br />
          See it happen.
        </h1>

        <p className="hero-description">
          Create interactive polls and watch responses appear
          instantly — no refresh required.
        </p>

        <div className="feature-list">
          <div className="feature">
            <div className="feature-icon">⚡</div>
            <div>
              <strong>Live results</strong>
              <span>Votes update instantly</span>
            </div>
          </div>

          <div className="feature">
            <div className="feature-icon">🔗</div>
            <div>
              <strong>Share anywhere</strong>
              <span>One link for your audience</span>
            </div>
          </div>

          <div className="feature">
            <div className="feature-icon">🛡️</div>
            <div>
              <strong>Secure polling</strong>
              <span>Protected poll management</span>
            </div>
          </div>
        </div>
      </div>

      <div className="hero-orb orb-one"></div>
      <div className="hero-orb orb-two"></div>
    </div>

    <div className="login-panel">
      <div className="login-box">
        <span className="badge">
          {mode === "login" ? "WELCOME BACK" : "GET STARTED"}
        </span>

        <h2>
          {mode === "login"
            ? "Welcome back"
            : "Create your account"}
        </h2>

        <p className="login-subtitle">
          {mode === "login"
            ? "Sign in to create and manage your polls."
            : "Start creating live polls in seconds."}
        </p>

        <form onSubmit={handleAuth}>
          {mode === "register" && (
            <div className="input-group">
              <label>Your name</label>
              <input
                placeholder="Enter your name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>
          )}

          <div className="input-group">
            <label>Email address</label>
            <input
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div className="input-group">
            <label>Password</label>
            <input
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <button
            className="primary-btn login-btn"
            disabled={loading}
          >
            {loading
              ? "Please wait..."
              : mode === "login"
              ? "Sign in →"
              : "Create account →"}
          </button>
        </form>

        <div className="login-switch">
          <span>
            {mode === "login"
              ? "New to LivePoll?"
              : "Already have an account?"}
          </span>

          <button
            className="text-btn"
            onClick={() =>
              setMode(mode === "login" ? "register" : "login")
            }
          >
            {mode === "login" ? "Create account" : "Sign in"}
          </button>
        </div>
      </div>
    </div>
  </section>
)}
{isLoggedIn && !pollId && (
  <section className="create-card">

    <span className="badge">CREATE A POLL</span>

    <h2>Create your live poll</h2>

    <p>
      Ask a question, add your options, and share the poll with your audience.
      Results will update in real time.
    </p>

    <form onSubmit={createPoll}>

      <div className="input-group">
        <label>Poll question</label>

        <input
          type="text"
          placeholder="e.g. What is your favorite programming language?"
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          required
        />
      </div>

      <div className="options-heading">
        <label>Answer options</label>
        <span>{options.length}/6</span>
      </div>

      {options.map((option, index) => (
        <div className="option-input-row" key={index}>

          <input
            type="text"
            placeholder={`Option ${index + 1}`}
            value={option}
            onChange={(e) =>
              updateOption(index, e.target.value)
            }
            required
          />

          {options.length > 2 && (
            <button
              type="button"
              className="remove-option"
              onClick={() => {
                setOptions(
                  options.filter((_, i) => i !== index)
                );
              }}
            >
              ×
            </button>
          )}

        </div>
      ))}

      {options.length < 6 && (
        <button
          type="button"
          className="add-option"
          onClick={addOption}
        >
          + Add option
        </button>
      )}

      <button
        type="submit"
        className="primary-btn create-poll-btn"
        disabled={loading}
      >
        {loading ? "Creating poll..." : "Create Live Poll →"}
      </button>

    </form>

  </section>
)}

        {poll && (
          <section className="poll-card">
            <span className="badge">LIVE RESULTS</span>

            <h1>{poll.question}</h1>

            <p className="live-status">
              <span className="live-dot"></span>
              Live voting · {totalVotes} vote{totalVotes === 1 ? "" : "s"}
            </p>

            <div className="options">
              {poll.options.map((option) => {
                const count = Number(option.count || 0);
                const percentage =
                  totalVotes > 0
                    ? Math.round((count / totalVotes) * 100)
                    : 0;

                return (
                  <button
                    className="option"
                    key={option.id}
                    onClick={() => vote(option.id)}
                  >
                    <div className="option-top">
                      <span>{option.text}</span>
                      <strong>{percentage}%</strong>
                    </div>

                    <div className="bar">
                      <div
                        className="bar-fill"
                        style={{ width: `${percentage}%` }}
                      ></div>
                    </div>

                    <small>
                      {count} vote{count === 1 ? "" : "s"}
                    </small>
                  </button>
                );
              })}
            </div>

            <div className="share-box">
              <span>Share this poll</span>
              <input
                readOnly
                value={window.location.href}
                onFocus={(e) => e.target.select()}
              />
              <button
                className="secondary-btn"
                onClick={() => navigator.clipboard.writeText(window.location.href)}
              >
                Copy link
              </button>
            </div>
          </section>
        )}
      </main>
    </div>
  );
}

export default App;