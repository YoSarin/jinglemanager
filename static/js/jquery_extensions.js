// Author:  Jacek Becela
// Source:  http://gist.github.com/399624
// License: MIT

jQuery.fn.multi_click = function(single_click_callback, double_click_callback, triple_click_callback, timeout) {
    return this.each(function(){
        var clicks = 0, self = this;
        jQuery(this).click(function(event){
            clicks++;
            if (clicks == 1) {
                setTimeout(function(){
                    if(clicks == 1) {
                        single_click_callback.call(self, event);
                    } else if (clicks == 2){
                        double_click_callback.call(self, event);
                    } else if (clicks == 3 && triple_click_callback){
                        triple_click_callback.call(self, event);
                    }
                    clicks = 0;
                }, timeout || 300);
            }
        });
    });
}
