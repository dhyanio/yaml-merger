# YAML Content Merger

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.23-61CFDD.svg?style=flat-square)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Use cases
- [Swagger](./examples/swagger/): Merge multiple swagger files into a swagger file, support JSON/YAML.

A simple yet powerful tool written in Go for merging multiple YAML files. This tool allows you to merge complex YAML structures such as lists and maps, with the ability to ignore specific keys during the merge process.

## 🔧 Features
- **Recursive Merging**: Handles nested structures, including lists and maps.
- **Ignore Keys**: Optionally exclude specific keys from merging.
- **Error-Resilient**: Provides clear error messages for parsing and merging issues.
- **Simple CLI Usage**: Merge files directly from the command line.
- **Dry Run**: Preview merged content without saving.
- **Output**: Output merged content to the specified file.
- **Custom Merge Strategies**: Choose between merge (deep merge) and override for key conflicts.

## 🚀 Installation
To install the tool, use the go install command:

```bash
go install github.com/dhyanio/yaml-merger@latest
```
This will install the yaml-merger binary, which you can run from your terminal.

## Prerequisites
Go 1.16+ is required for the go install command.

## 📚 Usage
After installing the tool, you can run it from the command line:

```bash
yaml-merger [options] <file1.yaml> <file2.yaml> [...]
```

## 🔍 Example:
```bash
yaml-merger config1.yaml config2.yaml config3.yaml
```
This will merge the contents of config1.yaml, config2.yaml, and config3.yaml and output the result to stdout.

## Options:
- `--ignore`: A comma-separated list of keys to ignore during the merge process.
  ```bash
  yaml-merger --ignore=key1,key2 file1.yaml file2.yaml
  ```
  This will merge file1.yaml and file2.yaml while ignoring the specified keys.

- `--dry-run`: Simulates the process without writing to a file. Instead, it prints the merged content.
  ```bash
  yaml-merger --dry-run file1.yaml file2.yaml
  ```
  This will display the merged YAML without saving it to a file.

- `--merge-strategy`: Output merged YAML content to the specified file. If not, it prints the output to stdout.
  ```bash
  yaml-merger --merge-strategy merge file1.yaml file2.yaml
  ```
  Deep Merge: Use the default "merge" strategy for recursive merging.

  ```bash
  yaml-merger --merge-strategy override file1.yaml file2.yaml
  ```
  Override: Use the "override" strategy to override values when keys clash.

- `--output`: Output merged YAML content to the specified file. If not, it prints the output to stdout.
  ```bash
  yaml-merger --output result.yaml file1.yaml file2.yaml
  ```
  Merge YAML files and write to file (without dry run)

### Example Files:
file1.yaml
```yaml
app:
  name: "MyApp"
  version: "1.0"
environment: "production"
features:
  - "logging"
  - "monitoring"
```
file2.yaml

```yaml
app:
  version: "2.0"
  author: "Dhyanio"
features:
  - "caching"
  - "alerts"
```
Merged Output:
```yaml
app:
  name: "MyApp"
  version: "2.0"
  author: "Dhyanio"
environment: "production"
features:
  - "Logging"
  - "Monitoring"
  - "caching"
  - "alerts"
```
## 🛠 Logging
The tool uses https://github.com/dhyanio/gogger for structured logging. Logs are set to INFO by default, but you can enable debug-level logging for detailed output.

### Enable Debug Logging:
```bash
LOG_LEVEL=debug yaml-merger file1.yaml file2.yaml
```

##❗Error Handling
The tool provides detailed error messages in case of:

- File Read Errors: The tool will display an error message if a file cannot be read.
- YAML Parse Errors: If there are syntax issues in the YAML files.
- Type Mismatch: If incompatible types are encountered during merging (e.g., list vs map).

## 🤝 Contributing
Contributions are welcome! Please open an issue or submit a pull request on GitHub.
1. Fork the repository.
2. Create your feature branch (git checkout -b feature/new-feature).
3. Commit your changes (git commit -am 'Add new feature').
4. Push to the branch (git push origin feature/new-feature).
5. Create a new Pull Request.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## ❤️ Acknowledgements

Thanks to the Go community for their support and contributions.
