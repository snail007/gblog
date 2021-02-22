;jQuery(function () {
    if(window.init){
        for(var i=0;i<window.init.length;i++){
            window.init[i]();
        }
    }
});