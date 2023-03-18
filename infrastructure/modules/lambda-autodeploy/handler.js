const { LambdaClient, GetFunctionCommand, UpdateFunctionCodeCommand } = require("@aws-sdk/client-lambda");

const lmbClient = new LambdaClient()

const tagS3Bucket = 'littleurl-autodeploy-s3-bucket'
const tagS3Key = 'littleurl-autodeploy-s3-key'
const tagEnabled = 'littleurl-autodeploy-enabled'

exports.handler = async (event, context) => {
  console.log('Received event:', JSON.stringify(event, null, 2));

  for (const record of event.Records) {
    // Get the object from the event
    const s3Bucket = record.s3.bucket.name
    const s3Key = decodeURIComponent(record.s3.object.key.replace(/\+/g, ' '))

    // check the object is a zip file
    if (!s3Key.endsWith('.zip')) {
      console.warn('Unrecognised file type', s3Key)
    }

    // get current function config
    const funcName = `${process.env.FUNCTION_NAME_PREFIX}${s3Key.replace('.zip', '')}`
    const func = await lmbClient.send(new GetFunctionCommand({
      FunctionName: funcName,
    }))

    // check autodeploy is enabled
    if (func.Tags[tagEnabled] !== 'true') {
      console.info(`Autodeploy disabled for function: ${func.FunctionName}`)
      continue
    }

    // check autodeploy tags
    if (func.Tags[tagS3Bucket] !== s3Bucket || func.Tags[tagS3Key] !== s3Key) {
      console.error('Mismatching autodeploy tags', {
        eventS3Bucket: s3Bucket,
        eventS3Key: s3Key,
        tagS3Bucket: func.Tags[tagS3Bucket],
        tagS3Key: func.Tags[tagS3Key],
      })
      continue
    }

    // update function
    console.info(`Updating function ${funcName}`)
    await lmbClient.send(new UpdateFunctionCodeCommand({
      FunctionName: funcName,
      S3Bucket: s3Bucket,
      S3Key: s3Key,
    }))
  }
}
