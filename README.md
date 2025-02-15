# ai-commit

Are you tired to think about a commit message at the end of making changes on your code?
Have you often `Update README.md` or `test` as message in your git log? Let create AI for you a meanful message based on last `git diff`. Here is a Git Commit Hook with AI Support

Look at the [demo](https://eumel8.github.io/ai-commit/recording.html) how to work on a Terraform manifest for OpenTelekomCloud

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


# customization

ai-commit works out of the bos with a local Ollama instance. Nevertheless you can customize such things in the `commit-msg` file:

* OLLAMA_API: API endpoint of OLLAMA instance
* MODEL: Name of the usesd LLM (available on the OLLAMA host, can be any LLM such llamam, mistral, deepseek-r1)
* PROMPT: exact task explained to the AI
* COMMIT_MSG: Lots of magic to modify the output and make a conform Git commit message
   
# super lazy mode

```
alias ac='git add . && git commit -m ""'
```

# Credits

Frank Kloeker f.kloeker@telekom.de

Life is for sharing. If you have an issue with the code or want to improve it, feel free to open an issue or an pull request.
