import random
from math import gcd

def get_primes_sieve(n: int) -> list[int]:
  prime = [True] * (n + 1)
  p = 2

  while p * p <= n:
    if prime[p]:
      for i in range(p * p, n + 1, p):
        prime[i] = False
    p += 1
  
  res = []
  for p in range(2, n + 1):
    if prime[p]:
      res.append(p)
  
  return res

def get_coprime(primes: list[int], phi: int) -> int:
  while len(primes) > 0:
    e = random.choice(primes)
    if gcd(e, phi) == 1:
      return e
    else:
      primes.remove(e)
  else:
    raise ValueError('No number coprime with phi was found')

def power_mod(base: int, exponent: int, mod: int) -> int:
  result = 1
  while exponent >= 1:
    if exponent % 2 == 1:
      result = (result * base) % mod
      exponent -= 1
    else:
      base = (base * base) % mod
      exponent //= 2
  return result

def encrypt_message(message: list[int], public_key: int, mod: int) -> list[int]:
  return [power_mod(char, public_key, mod) for char in message]

def decrypt_message(message: list[int], private_key: int, mod: int) -> list[int]:
  return [power_mod(char, private_key, mod) for char in message]

if __name__ == '__main__':
  primes = get_primes_sieve(100)
  high_primes = primes[len(primes) // 2:]
  p = random.choice(high_primes)
  q = random.choice(high_primes)
  n = p * q
  phi = (p - 1) * (q - 1)
  public_key = get_coprime(high_primes, phi)
  private_key = pow(public_key, -1, phi)

  print(f'Random values:')
  print(f'p\t{p}')
  print(f'q\t{q}')
  print(f'n\t{n}')
  print(f'phi\t{phi}')
  print(f'e\t{public_key}')
  print(f'd\t{private_key}')

  print('----------------')
  message = input('Insert message: ')
  unicode_characters = [ord(char) for char in message]
  print(f'UTF8:\t{unicode_characters}')
  encrypted_message = encrypt_message(unicode_characters, public_key, n)
  print(f'RSA:\t{encrypted_message}')
  decrypted_message = decrypt_message(encrypted_message, private_key, n)
  unicode_message = ''.join([chr(code) for code in decrypted_message])
  print(f'Decrypt:\t{unicode_message}')