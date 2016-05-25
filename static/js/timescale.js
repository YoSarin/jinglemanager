Timescale = function () {
    var self = this;
    this.items = [];

    this.start = function() {
        var start = null;
        $.each(self.items, function(k, v) {
            if (start == null || v.start < start) {
                start = v.start
            }
        });
        return start;
    }

    this.end = function() {
        var end = null;
        $.each(self.items, function(k, v) {
            var slotEnd = v.start.getTime() + v.duration * 1000
            if (end == null || (slotEnd) > end.getTime()) {
                end = new Date(slotEnd)
            }
        });
        return end;
    }

    this.duration = function() {
        return (self.end().getTime() - self.start().getTime()) / 1000 / 60
    }

    this.add = function(item) {
        self.items.push(item)
    }

    this.displayList = function() {
        var out = []
        l = self.items
        l.sort(function (a, b) {
            if (a.start < b.start) {
                return -1;
            } else if (a.start > b.start) {
                return 1;
            }
            return 0;
        });
        var lastEnd = null;
        $.each(l, function (k, v) {
            var gap = 0;
            if (lastEnd != null) {
                gap = (v.start.getTime() - lastEnd) / 1000 / 60;
            }
            lastEnd = (v.start.getTime() + v.duration * 1000);
            out.push({
                duration: v.duration / 60,
                start: v.start,
                end: new Date(v.start.getTime() + v.duration * 1000),
                gapBefore: gap
            });
        });
        return out;
    }

    return this
}

Slot = function (start, duration) {
    this.start = new Date(Date.parse(start));
    // go api returns it as a nanoseconds (facepalm)
    this.duration = duration/1000000000
    return this
}
