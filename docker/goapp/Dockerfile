FROM golang:1.18.4

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive

# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Install ImageMagick deps
RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev libmagickwand-dev redis --no-install-recommends

ENV IMAGEMAGICK_VERSION=7.1.0.43

RUN cd && \
	wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib

WORKDIR /app

COPY . /app

RUN export CGO_LDFLAGS="\
-Wl,-Bstatic \
    `pkg-config --libs MagickWand MagickCore` \
     -ljbig -ltiff -ljpeg -lwebp -llzma -lfftw3 -lbz2 -lgomp \
-Wl,-Bdynamic \
    -llcms2 -llqr-1 -lglib-2.0 -lpng12 -lxml2 -lz -lm -ldl \
"

RUN go install && \
    go build ./main.go

CMD /app/main