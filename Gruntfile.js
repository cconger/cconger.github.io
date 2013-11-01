module.exports = function(grunt) {

  grunt.initConfig({
    jade: {
      files: {
        expand: true,
        cwd: "src/",
        src: "**/*.jade",
        dest: "./",
        ext: ".html"
      }
    },
    stylus: {
      files: {
        expand: true,
        cwd: "src/scss",
        src: "**/*.styl",
        dest: "./css",
        ext: ".css"
      }
    },
    watch: {
      files: ['src/**/*'],
      tasks: ['default']
    }
  });

  grunt.loadNpmTasks('grunt-contrib-jade');
  grunt.loadNpmTasks('grunt-contrib-stylus');
  grunt.loadNpmTasks('grunt-contrib-watch');

  grunt.registerTask('default', ['jade', 'stylus']);
};

