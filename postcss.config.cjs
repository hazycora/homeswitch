const autoprefixer = require('autoprefixer')
const postcssNesting = require('postcss-nesting')
const atImport = require('postcss-import')

const config = {
	plugins: [atImport(), postcssNesting(), autoprefixer]
}

module.exports = config
