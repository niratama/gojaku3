/* global $, remark, TweetPanel, ProgressPanel  */

$(function () {
  'use strict';
  var slideshow = remark.create({
    ratio: '4:3',
    highlightLanguage: 'go',
    highlightStyle: 'github'
  });
  var progressPanel = new ProgressPanel($("#progress-panel"), {
    timer: {
    },
    page: {
      total: slideshow.getSlideCount(),
      current: slideshow.getCurrentSlideNo()
    }
  });
  var tweetPanel = new TweetPanel($('#tweet-panel'));

  slideshow.on('showSlide', progressPanel.showSlideHandler());

  $(window).on('resize', function () {
    tweetPanel.resize();
    progressPanel.resize();
  });

  var query = 'golang,go弱,go言語,#gojaku';
  var ws = new WebSocket('ws://localhost:3000/search?q=' +
    encodeURIComponent(query));
  ws.onmessage = function (e) {
    var msg = JSON.parse(e.data);
    tweetPanel.add(msg);
  };
  $(window).unload(function () {
    ws.close();
    ws = null;
  });
});

// vi:set sts=2 sw=2 et:
