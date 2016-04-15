package builder

// AquaLambdaURL is the URL to the Lambda installation file for Aqua
var AquaLambdaURL = "https://github.com/ArjenSchwarz/aqua/releases/download/latest/aqua_lambda.zip"

// Helloworld64 is a base64 encoded Hello World NodeJS app (zipfile)
var Helloworld64 = "UEsDBBQAAAAIAFqghkjgANcUYwAAAG4AAAAIABwAaW5kZXguanNVVAkAA8veBFfN3gRXdXgLAAEE9QEAAAQUAAAALYxBCsJAEATveUWTU4KyDzDkITnG3dYIZkZ2ZiVB/HsWsW4FRXF7aXYLyyzpyYwRtyLRHyod3xQ/I6o4N+/xaVD5a7ASI5m6dtICqyV8oRFzvpe1ql1anPB7hKumvR+a73AAUEsBAh4DFAAAAAgAWqCGSOAA1xRjAAAAbgAAAAgAGAAAAAAAAQAAAKSBAAAAAGluZGV4LmpzVVQFAAPL3gRXdXgLAAEE9QEAAAQUAAAAUEsFBgAAAAABAAEATgAAAKUAAAAAAA=="

// TrustDocument is a Lambda trustdocument for Role creation
var TrustDocument = `{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Principal": {"Service": "lambda.amazonaws.com"},
    "Action": "sts:AssumeRole"
  }
}
`

// BasicRole is the IAM Role configuration for a basic Lambda role
var BasicRole = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
`

// S3Role is the IAM Role configuration for a Lambda role with S3 access
var S3Role = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": [
        "arn:aws:s3:::*"
      ]
    }
  ]
}
`

// AquaRole is the IAM Role configuration required for Aqua to run
var AquaRole = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "apigateway:*",
        "lambda:AddPermission",
        "lambda:CreateFunction",
        "lambda:GetFunctionConfiguration"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "iam:GetRole",
        "iam:PassRole"
      ],
      "Resource": "arn:aws:iam::*:role/*"
    }
  ]
}
`
