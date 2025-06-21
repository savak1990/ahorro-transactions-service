terraform {
  backend "s3" {
    bucket         = "ahorro-app-state"
    key            = "pipeline/ahorro-transactions-service/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "ahorro-app-state-lock"
    encrypt        = true
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.0"
}

provider "aws" {
  region = "eu-west-1"

  default_tags {
    tags = {
      Project   = "ahorro-app"
      Service   = "ahorro-transactions-service-pipeline"
      Terraform = "true"
    }
  }
}

locals {
  project_name    = "ahorro-transactions-service"
  codebuild_name  = "${local.project_name}-build"
  github_owner    = "savak1990"
  github_repo     = "ahorro-transactions-service"
  github_branch   = "main"
  artifact_bucket = "ahorro-artifacts"
}

data "aws_codestarconnections_connection" "github" {
  name = "ahorro-app-github-connection"
}

resource "aws_codepipeline" "go_pipeline" {
  # Note that s3 bucket name comes from pipeline name, and can be truncated
  name     = "transactions"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = local.artifact_bucket
    type     = "S3"
  }

  stage {
    name = "Source"
    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["source_output"]
      configuration = {
        ConnectionArn    = data.aws_codestarconnections_connection.github.arn
        FullRepositoryId = "${local.github_owner}/${local.github_repo}"
        BranchName       = local.github_branch
        DetectChanges    = "true"
      }
    }
  }

  stage {
    name = "Build"
    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      version          = "1"
      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]
      configuration = {
        ProjectName = aws_codebuild_project.go_build.name
      }
    }
  }
}

resource "aws_codebuild_project" "go_build" {
  name         = local.codebuild_name
  service_role = aws_iam_role.codebuild_role.arn

  artifacts {
    type     = "S3"
    location = local.artifact_bucket
    # The path should match the name of the CodePipeline to be in the same bucket
    path                = "transactions"
    packaging           = "ZIP"
    name                = "ahorro-transactions-service.zip"
    encryption_disabled = true
  }

  environment {
    compute_type    = "BUILD_GENERAL1_SMALL"
    image           = "aws/codebuild/standard:7.0"
    type            = "LINUX_CONTAINER"
    privileged_mode = false
  }

  source {
    type                = "GITHUB"
    location            = "https://github.com/${local.github_owner}/${local.github_repo}.git"
    git_clone_depth     = 1
    buildspec           = "pipeline/buildspec.yml"
    report_build_status = true

  }

  logs_config {
    cloudwatch_logs {
      status     = "ENABLED"
      group_name = "/aws/codebuild/${local.project_name}-build"
    }
  }

  project_visibility = "PUBLIC_READ"
}

resource "aws_iam_role" "codebuild_role" {
  name = "${local.project_name}-codebuild-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Principal = {
        Service = "codebuild.amazonaws.com"
      },
      Effect = "Allow",
    }]
  })
}

resource "aws_iam_role_policy_attachment" "codebuild_policy" {
  role       = aws_iam_role.codebuild_role.name
  policy_arn = "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess"
}

resource "aws_iam_policy" "s3_artifacts" {
  name = "${local.project_name}-s3-artifacts-access"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "s3:GetObjectVersion",
        "s3:GetBucketAcl",
        "s3:GetBucketLocation"
      ],
      Resource = [
        "arn:aws:s3:::${local.artifact_bucket}",
        "arn:aws:s3:::${local.artifact_bucket}/*"
      ]
    }]
  })
}

resource "aws_iam_role_policy_attachment" "codebuild_s3" {
  role       = aws_iam_role.codebuild_role.name
  policy_arn = aws_iam_policy.s3_artifacts.arn
}

resource "aws_iam_policy" "codebuild_logs" {
  name = "${local.project_name}-codebuild-logs"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "codebuild_logs_attach" {
  role       = aws_iam_role.codebuild_role.name
  policy_arn = aws_iam_policy.codebuild_logs.arn
}

resource "aws_iam_role" "codepipeline_role" {
  name = "${local.project_name}-codepipeline-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect    = "Allow",
      Principal = { Service = "codepipeline.amazonaws.com" },
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "codepipeline_logs" {
  name = "${local.project_name}-codepipeline-logs"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "codepipeline_logs_attach" {
  role       = aws_iam_role.codebuild_role.name
  policy_arn = aws_iam_policy.codepipeline_logs.arn
}

resource "aws_iam_role_policy_attachment" "codepipeline_codebuild" {
  role       = aws_iam_role.codepipeline_role.name
  policy_arn = "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess"
}

resource "aws_iam_role_policy_attachment" "codepipeline_basic" {
  role       = aws_iam_role.codepipeline_role.name
  policy_arn = "arn:aws:iam::aws:policy/AWSCodePipeline_FullAccess"
}

resource "aws_iam_role_policy_attachment" "codepipeline_s3_custom" {
  role       = aws_iam_role.codepipeline_role.name
  policy_arn = aws_iam_policy.s3_artifacts.arn
}

resource "aws_iam_policy" "codepipeline_pass_codebuild_role" {
  name = "${local.project_name}-codepipeline-pass-codebuild-role"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect   = "Allow",
      Action   = "iam:PassRole",
      Resource = aws_iam_role.codebuild_role.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "codepipeline_pass_codebuild_role" {
  role       = aws_iam_role.codepipeline_role.name
  policy_arn = aws_iam_policy.codepipeline_pass_codebuild_role.arn
}

resource "aws_iam_policy" "codepipeline_codestar_connections" {
  name = "${local.project_name}-codepipeline-codestar-connections"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Action = [
        "codestar-connections:UseConnection"
      ],
      Resource = data.aws_codestarconnections_connection.github.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "codepipeline_codestar_connections" {
  role       = aws_iam_role.codepipeline_role.name
  policy_arn = aws_iam_policy.codepipeline_codestar_connections.arn
}
