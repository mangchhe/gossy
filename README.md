# Gossy

**Gossy** is a versatile AWS CLI tool designed for efficient management across AWS services. It simplifies connecting to EC2 instances via SSM and setting up port forwarding to RDS instances.

## Features

- **AWS Profile Management**: Easily select and switch between AWS profiles
- **EC2 Instance Connection**: Connect to EC2 instances using AWS Systems Manager (SSM)
- **RDS Port Forwarding**: Set up port forwarding to RDS instances through EC2 instances
- **Session Management**: Remember your last session for quick reconnection
- **Interactive CLI**: User-friendly interface with interactive prompts

## Installation

Install via Homebrew:

```sh
brew tap mangchhe/gossy
brew install gossy
```

## Prerequisites

- AWS CLI installed and configured with profiles
- AWS Systems Manager (SSM) agent installed on your EC2 instances
- Appropriate IAM permissions for SSM and RDS access

## AWS Profile Setup

To use Gossy, you need to have AWS profiles configured. Here's how to set them up:

1. Install the AWS CLI if you haven't already:
   ```sh
   brew install awscli
   ```

2. Configure your AWS profiles:
   ```sh
   aws configure --profile myprofile
   ```
   You'll be prompted to enter:
   - AWS Access Key ID
   - AWS Secret Access Key
   - Default region name (e.g., ap-northeast-2)
   - Default output format (e.g., json)

3. You can set up multiple profiles by repeating the command with different profile names:
   ```sh
   aws configure --profile another-profile
   ```

4. Verify your profiles:
   ```sh
   aws configure list-profiles
   ```

Your profiles will be stored in `~/.aws/credentials` and `~/.aws/config` files.

## Usage

### Basic Usage

Simply run the command:

```sh
gossy
```

This will:
1. Show a list of your AWS profiles
2. Allow you to select an EC2 instance
3. Choose between connecting directly to the instance or setting up port forwarding to an RDS instance

<details>
<summary><strong>Example Output (Click to expand)</strong></summary>

```
~ gossy
? Select an AWS Profile:  [Use arrows to move, type to filter]
> Connect with Last Session
  development
  production
  staging
  testing
  personal
  default

~ gossy
? Select an AWS Profile: development
? Select an EC2 Instance:  [Use arrows to move, type to filter]
> app-server-1   (i-0a1b2c3d4e5f67890) - running
  app-server-2   (i-0123456789abcdef0) - running
  db-bastion     (i-abcdef0123456789) - running

? Select connection type:  [Use arrows to move, type to filter]
> Instance (SSM)
  DB (Port Forwarding)
```
</details>

### Quick Reconnect

To quickly reconnect to your last session:

```sh
gossy
```

Then select "Connect with Last Session" from the profile menu.

### Profile Selection

To explicitly select a profile:

```sh
gossy profile
```

This command skips the "Connect with Last Session" option and shows all your AWS profiles alphabetically.

## How It Works

Gossy uses:
- AWS SDK for Go to interact with AWS services
- AWS SSM to establish secure connections to EC2 instances
- Port forwarding through SSM to connect to RDS instances
- Local storage to remember your last session

## Development

### Requirements

- Go 1.16 or higher
- AWS SDK for Go v2

### Building from Source

```sh
git clone https://github.com/mangchhe/gossy.git
cd gossy
go build
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.