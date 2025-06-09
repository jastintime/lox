.POSIX:
.PHONY: clean all 

CC = cc
CFLAGS = -O2 -g -Wall -Wextra -Winit-self -Wuninitialized -pedantic -Wunreachable-code
BUILD_DIR = build
OBJS = $(BUILD_DIR)/main.o $(BUILD_DIR)/chunk.o $(BUILD_DIR)/memory.o $(BUILD_DIR)/debug.o $(BUILD_DIR)/value.o

all: lox

lox: $(OBJS)
	$(CC) $(CFLAGS) -o $@ $(OBJS)

$(BUILD_DIR)/chunk.o: chunk.c chunk.h memory.h common.h value.h
	mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

$(BUILD_DIR)/memory.o: memory.c memory.h
	mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

$(BUILD_DIR)/main.o: main.c common.h chunk.h debug.h
	mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

$(BUILD_DIR)/debug.o: debug.c debug.h 
	mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

$(BUILD_DIR)/value.o: value.c memory.h value.h
	mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

clean:
	-rm -f $(OBJS) lox
	-rmdir $(BUILD_DIR)
