
CC=gcc
TARGET=./foo
SOURCE=./foo.c

run : ${TARGET}
	${TARGET}

clean :
	-rm ${TARGET}

${TARGET} : ${SOURCE}
	${CC} $^ -o $@

.PHONY : run clean
