//当前页面是哪一页 ，对应的页码变色
$(document).ready(function () {
    //如果页数大于5  第六个以后就显示成。。。 然后再显示最后两页
    url_controller = "index"
    url_key_word = ""
    //alert(page_type)
    if (page_type == "search") {
        //  alert("进来了")
        url_controller = "search_article"
        url_key_word = "/?key_word=" + key_word + ""
    }
    // alert(url_controller)
    $("#index_article_page").html(function () {
            var a_h = ""
            //上一页
            if (currentPage == 1) {
                a_h = "<li  class=\"previous disabled\"><span><span aria-hidden=\"true\">&laquo;</span></span></li>"
            } else {
                a_h = "<li><a href=\"/" + url_controller + "/" + (currentPage - 1) + "" + url_key_word + "\" style=\"cursor:pointer;\" aria-label=\"Previous\"><span aria-hidden=\"true\">&laquo;</span></a></li>"
            }
            for (var i = 0; i < pageCount; i++) {
                var c = i + 1
                var a_active = "<li class = \"active\"><a class =\"user_article_page_a\" style=\"cursor:pointer;\" href=\"/" + url_controller + "/" + c + "" + url_key_word + "\">" + c + "</a></li>"
                var a_previous = "<li class=\"previous disabled\"><span>...</span></li>"
                var a_normal = "<li><a class =\"user_article_page_a\" style=\"cursor:pointer;\" href=\"/" + url_controller + "/" + c + "" + url_key_word + "\">" + c + "</a></li>"
                //如果页码是当前页码  就设置高亮 li class = "active"
                if (c == currentPage) {
                    a_h = a_h + a_active
                } else {
                    //如何总页码大于5   五以后显示 。。。  然后再显示最后两页
                    if (pageCount > 5) {
                        //如果当前页面小于4  就显示 << 1 2 3 4 5 。。。 13 14 >>
                        if (currentPage < 4) {
                            if (c == pageCount || c == pageCount - 1) {
                                a_h = a_h + a_normal
                            } else if (c == 6) {
                                a_h = a_h + a_previous
                            } else if (c <= 5) {
                                a_h = a_h + a_normal
                            }
                            //如果当前页面大于等于4 小于7  就显示      << 1 2 3 4 5 6 。。。 13 14 >>
                        } else if (currentPage >= 4 && currentPage < 7) {
                            if (c == pageCount || c == pageCount - 1) {
                                a_h = a_h + a_normal
                            } else if (c == currentPage + 3) {
                                a_h = a_h + a_previous
                            } else if (c <= currentPage + 2) {
                                a_h = a_h + a_normal
                            }
                            //如果当前页面大于7  显示 << 1 2 。。5 6 7 8 9 。。。 13 14 >>
                        } else if (currentPage >= 7) {
                            if (c == pageCount || c == pageCount - 1) {
                                a_h = a_h + a_normal
                            } else if (c == 1 || c == 2) {
                                a_h = a_h + a_normal
                            } else if (c == currentPage + 3) {
                                a_h = a_h + a_previous
                            } else if (c == currentPage - 3) {
                                a_h = a_h + a_previous
                            } else if (c <= currentPage + 2 && c >= currentPage - 2) {
                                a_h = a_h + a_normal
                            }
                        }
                    } else {
                        a_h = a_h + a_normal
                    }
                }
            }
            //下一页
            if (currentPage == pageCount) {
                a_h = a_h + "<li class=\"previous disabled\"><span><span aria-hidden=\"true\">&raquo;</span></span></li>"
            } else {
                a_h = a_h + "<li><a href=\"/" + url_controller + "/" + (currentPage + 1) + "" + url_key_word + "\" style=\"cursor:pointer;\" aria-label=\"Next\"><span aria-hidden=\"true\">&raquo;</span></a></li>"
            }
            return a_h
        }
    );
});