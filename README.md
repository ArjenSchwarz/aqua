Aqua is a tool for quickly creating API Gateways for your AWS Lambda functions, that can run as a Lambda function.

# What does it do?

Aqua helps you to quickly create a hassle-free API Gateway for a Lambda function. What it creates is a very simple Gateway, that listens to POST requests and passes the form parameters on to the Lambda function. For what I need this is generally enough, and you can always change it afterwards to suit your needs better.

If you haven't created a Lambda function yet, you can provide this as well. Aqua is also capable of running as a Lambda function itself, and comes with a built-in shortcut for the installation.

In order to support this, Aqua can also show you your available roles and create several basic ones.

# What does it not do?

Aqua is *not* a tool for managing your API Gateways or Lambda functions. If you wish to use something that does that there are plenty of other options available. This one is designed only for simple functionalities.

# Usage

An explanation of all options is available through Aqua's help.

```bash
$ aqua --help
Usage:
  aqua [flags]
  aqua [command]

Available Commands:
  install     Install Aqua as a Lambda function
  role        Display or create roles

Flags:
  -k, --apikey                  Endpoint can only be accessed with an API key
  -a, --authentication string   The Authentication method to be used (default "NONE")
  -f, --file string             The zip file for your Lambda function, either locally or http(s). The file will first be downloaded locally.
      --json                    Set to true to print output in JSON format
  -n, --name string             The name of the Lambda function
      --region string           The region for the lambda function and API Gateway (default "us-east-1")
  -r, --role string             The name of the IAM Role)
      --runtime string          The runtime of the Lambda function. (default "nodejs4.3")

Use "aqua [command] --help" for more information about a command.
```

A couple of examples, the output for all is the same as shown in the first example.

Create a Gateway for an existing Lambda function:

```bash
$ aqua --name existingFunction
Your endpoint is available at https://api4id.execute-api.us-east-1.amazonaws.com/prod/existingfunction
```

Create a Lambda function with sample code, complete with Gateway:

```bash
$ aqua --name newFunction --role roleName
```

Create a Lambda function with your own code, complete with Gateway:

```bash
$ aqua --name newFunction --role roleName --file path/to/file.zip
$ aqua --name newFunction --role roleName --file https://web/address/of/file.zip
```

## As Lambda function

If installed as a Lambda function, Aqua is capable only of adding a Gateway to a function or creating a Lambda function with sample code. You cannot give it code to install.

You will have to post the values, where the form fields have the same name as the flags when running it from the command line.

# Installation

Simply download the [latest release][latest] for your platform, and you can use it. You can place it somewhere in your PATH to ensure you can run it from anywhere.

The AWS configuration is read from the standard locations:

* Your environment variables (`AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY`).
* The values in your `~/.aws/credentials` file.
* Permissions from the IAM role the application has access to (when running on AWS)

[latest]: https://github.com/ArjenSchwarz/aqua/releases

## Installation to Lambda

First [download Aqua][latest] to your local machine, and then ensure you have a role with enough permissions. You can create that role using Aqua with the following command:

```bash
$ aqua role create --role RoleName --type aqua
```

Or you can manually create the role, with the permissions as shown [here][permissionslink].

You can then install Aqua using:

```bash
$ aqua install --role RoleName --name Aqua
```

This will download the latest version of the Lambda function, and install it with the name and role you specified. Other flags (region etc.) are available as well.

For security reasons, `aqua install` enforces the use of API keys. This means that after the installation you will need to assign those keys or set up a different authentication method. As Aqua can create unprotected endpoints for your Lambda functions, it is recommended you always require some form of authentication.

[permissionslink]: https://github.com/ArjenSchwarz/aqua

# Development

Aqua is deliberately limited in what it can do, but that doesn't mean it can't be improved. Work is ongoing, and will likely involve some under the hood restructuring.

If you wish to contribute you can always create Issues or Pull Requests. For Pull Request, just follow the standard pattern.

1. Fork the repository
2. Make your changes
3. Make a pull request that explains what it does

## TODO

* Create the releases on GitHub using Wercker
* Internal restructuring to make more sense (move stuff out of the root command)
* Add tests
* More documentation
