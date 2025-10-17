# Dougless Permissions Configuration Guide

## Overview

Dougless uses a **config-first permission model** with `.douglessrc` configuration files. This provides a cleaner, more maintainable approach compared to CLI flags, and allows permissions to be version-controlled with your project code.

In development, an interactive **two-prompt flow** helps you build `.douglessrc` incrementally: first you grant access, then you choose whether to persist it to config.

## Quick Start

### 1. Create a `.douglessrc` file

In your project directory, create a `.douglessrc` file with JSON configuration:

```json
{
  "permissions": {
    "read": ["./examples", "/tmp"],
    "write": ["/tmp", "./output"],
    "net": ["api.example.com", "localhost:3000"],
    "env": ["API_KEY"],
    "run": ["git"]
  }
}
```

### 2. Run your script

```bash
./dougless script.js
# Automatically discovers and loads .douglessrc from the script directory
```

That's it! No CLI flags needed.

## Permission Types

### `read` - File System Read Access

Controls which files and directories your script can read.

```json
{
  "permissions": {
    "read": [
      "./data",           // Read files in ./data directory
      "/tmp",             // Read files in /tmp
      "/home/user/config.json"  // Read specific file
    ]
  }
}
```

**Example Usage:**
```javascript
// This will work if ./data is in the read permissions
const content = await files.read('./data/file.txt');
```

### `write` - File System Write Access

Controls which files and directories your script can write to.

```json
{
  "permissions": {
    "write": [
      "./output",         // Write to ./output directory
      "/tmp",             // Write to /tmp
      "./logs/"           // Write to logs directory (note trailing slash)
    ]
  }
}
```

**Example Usage:**
```javascript
// This will work if ./output is in the write permissions
await files.write('./output/result.txt', 'Hello World');
```

### `net` - Network Access

Controls which hosts and ports your script can connect to.

```json
{
  "permissions": {
    "net": [
      "api.example.com",           // Allow HTTP/HTTPS to api.example.com
      "localhost:3000",            // Allow localhost on port 3000
      "*.github.com",              // Wildcard subdomain support
      "192.168.1.100:8080"        // IP address with port
    ]
  }
}
```

**Example Usage:**
```javascript
// This will work if api.example.com is in the net permissions
http.get('https://api.example.com/data', (err, response) => {
  console.log(response.body);
});
```

### `env` - Environment Variable Access

Controls which environment variables your script can read (future feature).

```json
{
  "permissions": {
    "env": [
      "API_KEY",
      "DATABASE_URL",
      "NODE_ENV"
    ]
  }
}
```

### `run` - Subprocess Execution

Controls which external programs your script can execute (future feature).

```json
{
  "permissions": {
    "run": [
      "git",
      "npm",
      "curl"
    ]
  }
}
```

## Interactive Two-Prompt Flow

When a script needs a permission that isn't granted yet (and you're in an interactive terminal), Dougless will prompt:

```
⚠️  Permission request: read access to './data/config.json'
Allow? (y/n): y
Save to .douglessrc? (y/n): y
✓ Granted and saved to .douglessrc
```

- If you choose not to save, the grant applies to the current session only
- In non-interactive environments (CI, pipes), prompts are disabled and access is denied unless configured

## Configuration Discovery

Dougless automatically searches for `.douglessrc` files starting from the script's directory:

1. **Script directory** - Looks for `.douglessrc` in the same directory as your script
2. **Current working directory** - Falls back to CWD if not found in script directory

## Path Patterns

### Relative Paths

Paths starting with `./` are relative to the script directory:

```json
{
  "permissions": {
    "read": ["./data", "./config"]
  }
}
```

### Absolute Paths

Paths starting with `/` are absolute:

```json
{
  "permissions": {
    "read": ["/tmp", "/var/log"]
  }
}
```

### Directory vs File

- **Directories**: Grant access to all files within the directory
  ```json
  "read": ["./data"]  // Grants access to ./data/file1.txt, ./data/file2.txt, etc.
  ```
  
- **Specific files**: Grant access to only that file
  ```json
  "read": ["./config.json"]  // Grants access only to ./config.json
  ```

## Best Practices

### 1. Version Control Your `.douglessrc`

Commit `.douglessrc` to your repository so team members and CI/CD systems use the same permissions:

```bash
git add .douglessrc
git commit -m "Add permissions config"
```

### 2. Use Minimal Permissions

Only grant the permissions your script actually needs:

```json
{
  "permissions": {
    "read": ["./data"],    // Only ./data, not the entire filesystem
    "write": ["./output"]  // Only ./output for results
  }
}
```

### 3. Document Permission Requirements

Add comments in your README explaining why permissions are needed:

```markdown
## Permissions

This script requires:
- `read: ["./data"]` - To read input configuration files
- `write: ["./output"]` - To write processed results
- `net: ["api.example.com"]` - To fetch external data
```

### 4. Different Configs for Different Environments

You can use different config files for development vs production:

```bash
# Development
cp .douglessrc.dev .douglessrc

# Production
cp .douglessrc.prod .douglessrc
```

## Why Config-First?

Dougless uses configuration files instead of CLI flags for permissions. This approach provides several advantages:

### CLI Flags Approach (Not Recommended)

```bash
# Every time you run, you need to remember all the flags
runtime --allow-read=/data --allow-write=/tmp --allow-net=api.example.com script.js
```

**Problems:**
- Hard to remember and type
- Not version controlled
- Can't be shared with team
- Error-prone

### Config-First Approach (Dougless)

```bash
# Create .douglessrc once
cat > .douglessrc << EOF
{
  "permissions": {
    "read": ["/data"],
    "write": ["/tmp"],
    "net": ["api.example.com"]
  }
}
EOF

# Then just run - permissions are automatic
./dougless script.js
```

**Advantages:**
- Simple, clean command
- Version controlled with your code
- Team shares same permissions
- Self-documenting
- Maintainable

## Troubleshooting

### "Permission denied" errors

If you get a permission denied error:

1. Check your `.douglessrc` file exists in the correct location
2. Verify the path/host is listed in the appropriate permission array
3. Check for typos in paths or hostnames
4. Ensure JSON is valid (use a JSON validator)

### Config not found

```
Error: no .douglessrc found in /path/to/script
```

**Solution**: Create a `.douglessrc` file in your script's directory.

### Invalid JSON

```
Error: failed to parse config: invalid character '}' looking for beginning of object key string
```

**Solution**: Validate your JSON syntax. Common issues:
- Trailing commas (not allowed in JSON)
- Missing quotes around keys or values
- Unmatched brackets or braces

## Example Configurations

### Simple Web Scraper

```json
{
  "permissions": {
    "read": ["./config.json"],
    "write": ["./output"],
    "net": ["api.example.com", "cdn.example.com"]
  }
}
```

### File Processor

```json
{
  "permissions": {
    "read": ["./input", "./templates"],
    "write": ["./output", "./logs"]
  }
}
```

### Full-Stack Development

```json
{
  "permissions": {
    "read": ["./src", "./config", "./public"],
    "write": ["./dist", "./logs", "/tmp"],
    "net": ["localhost:3000", "localhost:5432", "api.staging.example.com"]
  }
}
```

## Future Enhancements

Planned features for the config system:

- **Cascading configs** - Global `~/.douglessrc` + project `.douglessrc`
- **Comments support** - JSONC format for adding comments to config files
- **Config validation command** - `dougless config validate` to check your config
- **Config init command** - `dougless config init` to generate a template

## See Also

- [Permissions API](./permissions_api.md) - Programmatic permission checking (future)
- [Security Best Practices](./security.md) - Security guidelines (future)
- [CLI Reference](./cli.md) - Command-line options (future)

---

*Last Updated: October 17, 2025*
