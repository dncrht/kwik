var isControlPressed = false;

function resizeToc() {
  $('.js-toc').height($(window).height() - $('.js-toc').position().top - 106);
}

$(document).ready(function() {
  $('.js-toc').nestedToc({container: '.js-maincontent'});

  if ($('.js-toc').html() == '') {
    $('.js-toc').remove();
  }

  $('.js-sidepanel').css('position', 'fixed');
  if ($('.js-toc').length == 1) {
    $('.js-toc').css('overflow-y', 'auto');

    $(window).resize(function(){
      resizeToc();
    });

    resizeToc();
  }

  $(document).keyup(
    function(){
      isControlPressed = false;
    }
  ).keydown(
    function(e) {
      if (e.which == 17) {
        isControlPressed = true;
      }

      // Control + S
      if (e.which == 83 && isControlPressed) {
        $('.js-save').click();
        return false;
      }

      // Control + K
      if (e.which == 75 && isControlPressed) {
        $('.js-terms').focus();
        return false;
      }
    }
  );
});
