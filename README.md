# lumber

## Deploy image

* Ideal deploy flow
    1. git push origin master
        * CI at CircleCI
        * docker build at CircleCI
        * docker push to dockerhub
    2. k8s deploy with dockerhub latest images by Spinnaker

refs: http://tech.mercari.com/entry/2017/08/21/092743

## Architecture reference:

refs: https://speakerdeck.com/mercari/ja-golang-package-composition-for-web-application-the-case-of-mercari-kauru
