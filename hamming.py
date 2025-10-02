import random

def change_char_at_index(string: str, char: str, index: int) -> str:
  return string[:index] + char + string[index + 1:]

def calculate_redundunt_bits(length: int) -> int:
  for index in range(length):
    if(2 ** index >= length + index + 1):
      return index

def convert_to_hamming(data: str) -> str:
  redundunt_bits_count = calculate_redundunt_bits(len(data))
  result_data = ''
  data_pointer = 0
  exponent = 0

  for index in range(1, len(data) + redundunt_bits_count + 1):
    if index == 2 ** exponent:
      result_data = result_data + '0'
      exponent += 1
    else:
      result_data = result_data + data[data_pointer]
      data_pointer += 1

  for exponent in range(redundunt_bits_count):
    parity_count = 0
    count = 0
    index = 2 ** exponent

    while index <= len(result_data):
      count += 1
      if result_data[index - 1] == '1':
        parity_count += 1
      if count == 2 ** exponent:
        index += count
        count = 0
      index += 1
    
    result_data = change_char_at_index(result_data, str(parity_count % 2), 2 ** exponent - 1)

  return result_data

def decode_hamming(message: str) -> str:
  redundunt_bits_count = 0
  while 2 ** redundunt_bits_count < len(message) + 1:
    redundunt_bits_count += 1

  error_position = 0
  for exponent in range(redundunt_bits_count):
    parity_count = 0
    count = 0
    index = 2 ** exponent

    while index <= len(message):
      count += 1
      if message[index - 1] == '1':
        parity_count += 1
      if count == 2 ** exponent:
        index += count
        count = 0
      index += 1

    if parity_count % 2 != 0:
      error_position += 2 ** exponent

  if error_position != 0:
    error_index = error_position - 1
    print('!!! FOUND ERROR !!!')
    print(f'Error at index:\t{error_index}\n')
    flipped_bit = '0' if message[error_index] == '1' else '1'
    message = change_char_at_index(message, flipped_bit, error_index)

  decoded_message = ''
  exponent = 0
  for index in range(1, len(message) + 1):
    if index == 2 ** exponent:
      exponent += 1
    else:
      decoded_message += message[index - 1]
  
  return decoded_message

if __name__ == '__main__':
  message = '1011001'
  coded_message = convert_to_hamming(message)

  random_index = random.randint(0, len(coded_message) - 1)
  flipped_bit = '0' if coded_message[random_index] == '1' else '1'
  error_message = change_char_at_index(coded_message, flipped_bit, random_index)
  
  print(f'Coded message:\t{coded_message}')
  print(f'Error message:\t{error_message}')
  print(f'Index of error:\t{random_index}')
  print('\nDecoding message...\n')

  decoded_message = decode_hamming(error_message)

  print(f'Initial message:\t{message}')
  print(f'Decoded message:\t{decoded_message}')