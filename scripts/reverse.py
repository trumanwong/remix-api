from PIL import Image,ImageSequence
import sys


def convert(src, dst):
    with Image.open(src) as im:
        if im.is_animated:
            imgs = [frame.copy() for frame in ImageSequence.Iterator(im)]
            imgs.reverse()
            imgs[0].save(dst, save_all=True, append_images=imgs[1:], loop=0)


convert(sys.argv[1], sys.argv[2])
