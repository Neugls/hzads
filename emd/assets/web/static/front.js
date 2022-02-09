(function(){
    function isFullScreen() {
        return  !! (
            document.fullscreen || 
            document.mozFullScreen ||                         
            document.webkitIsFullScreen ||       
            document.webkitFullScreen || 
            document.msFullScreen 
        );
    }

    function fullscreen(){
            let el = document.documentElement;
            let rfs = el.requestFullScreen || el.webkitRequestFullScreen || el.mozRequestFullScreen || el.msRequestFullscreen;
            if (typeof rfs != "undefined" && rfs) {
                rfs.call(el);
            };
            return;
    }

    function loadAdsData() {
        return new Promise((resolve, reject) => {
            $.get("/ads").done(function(data){
                let ads = data.data;
                resolve(ads);
            }).fail(function(err){
                reject(err);
                console.log("load ads data error " + err);
            });
        });
    }

    function setupAds(ads){
        $("#ads-container").empty();
        


        let ads_html = "";
        ads.forEach(function(ad){
            ads_html += '<div class="ad-item" aid="'+ad.id+'">';
            if (ad.type == "image"){
                ads_html += '<img src="' + ad.content + '" />';
            }
            else if (ad.type == "video"){
                ads_html += '<video src="' + ad.content + '"  autoplay></video>';
            }else if (ad.type == "webpage"){
                ads_html += '<iframe src="' + ad.content + '"  frameborder="0" allowfullscreen></iframe>';
            }
            ads_html += '</div>';
        });
        $("#ads-container").html(ads_html);

        $("#ads-container").css("display", "block");
    }

    let timer = 0;
    let transition_time = 1000;
    let ads_interval = 5000;

    function playAds(ads){
        
        if (timer > 0 ) clearInterval(timer);
        if(ads.length > 1){
            let current_ad = 0;


            $("#ads-container .ad-item").eq(current_ad).css("left", "0");
            timer = setInterval( () => {
                let last = current_ad;
                current_ad = (current_ad + 1) % ads.length;
                $("#ads-container .ad-item").eq(last).animate({left: "-100vw"}, transition_time, ()=>{
                    $("#ads-container .ad-item").eq(last).css("left", "100vw");
                });
                $("#ads-container .ad-item").eq(current_ad).animate({left: "0"}, transition_time);
            }, ads_interval);
        }
    }

    let reloadInterval = 10*1000;

    function startShow(){
       const action = ()=>{
           loadAdsData().then(function(data){
                var ads = data;
                if( ads && ads.length > 0){
                    setupAds(ads);
                    playAds(ads);
                }
            });
        };
        

        setInterval(action, reloadInterval);
        action();
    }
    

    //开始启动
    $(document).on("click", 'a[href="#startShow"]', function(e){
        e.preventDefault();
        $("header").remove();
        startShow();
        fullscreen();
    });
})()