# Design memo

## Concept

* Contents management can be at git repository
* Site generator?
    * Unnecessary
    * Implement converter from markdown to html
        * Require minimum blog content: markdown files
        * Therefore, creating blog content by markdown files repository
    * Ref: [https://github.com/gohugoio/hugo](https://github.com/gohugoio/hugo)
* And need publish content server
    * Related: viewing entry

## Must features

* Viewing entry via Frontend app
* Posting markdown file

## Better features

* Preview draft entry
* Image upload to S3, or alternative uploadable storage
    * Write specific syntax for the static file as a relative link
    * Convert from specific syntax to public image storage URL
* CMS
* Entry holds status for publish/private
* Write comments
    * Consider how to recognize users

## Language, Middleware

* Backend
    * go, mysql
* Frontend
    * React
