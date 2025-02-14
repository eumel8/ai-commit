# ai-commit

Are you tired to think about a commit message at the end of making changes on your code?
Have you often `Update README.md` or `test` as message in your git log? Let create AI for you a meanful message based on last `git diff`. Here is a Git Commit Hook with AI Support

# preparation

* A running [Ollama](https://ollama.com/) instance running with API access, probably on your computer on `127.0.0.1:11434`
* Loaded LLM to use, e.g. [deepseek-coder-v2-lite-base-q4_k_m-gguf:latest](https://ollama.com/networkjohnny/deepseek-coder-v2-lite-base-q4_k_m-gguf), a DeepSeek-Coder-V2-Lite-Base with only 2GB

  ```bash
  # ollama pull networkjohnny/deepseek-coder-v2-lite-base-q4_k_m-gguf:latest
  # ollama serve
  ```

* required tools on your computer: curl, jq, sed
* copy `commit-msg` to `.git/hooks/commit-msg` in your Git repo and make it executable (`chmod +x`)
* make changes and review AI's work. If you are unhappy, just edit the message with `git commit --amend` bevore push
  
# super lazy mode

```
alias ac='git add . && git commit -m ""'
```
