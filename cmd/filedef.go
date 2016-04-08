package cmd

var helloworld64 = "UEsDBBQAAAAIAFqghkjgANcUYwAAAG4AAAAIABwAaW5kZXguanNVVAkAA8veBFfN3gRXdXgLAAEE9QEAAAQUAAAALYxBCsJAEATveUWTU4KyDzDkITnG3dYIZkZ2ZiVB/HsWsW4FRXF7aXYLyyzpyYwRtyLRHyod3xQ/I6o4N+/xaVD5a7ASI5m6dtICqyV8oRFzvpe1ql1anPB7hKumvR+a73AAUEsBAh4DFAAAAAgAWqCGSOAA1xRjAAAAbgAAAAgAGAAAAAAAAQAAAKSBAAAAAGluZGV4LmpzVVQFAAPL3gRXdXgLAAEE9QEAAAQUAAAAUEsFBgAAAAABAAEATgAAAKUAAAAAAA=="

var trustDocument = `{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Principal": {"Service": "lambda.amazonaws.com"},
    "Action": "sts:AssumeRole"
  }
}
`

var basicRole = `{
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

var s3Role = `{
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

var aquaRole = `{
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
