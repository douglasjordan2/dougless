package modules

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/dop251/goja"
	"github.com/google/uuid"
)

type Crypto struct {
  vm *goja.Runtime
}

func NewCrypto() *Crypto {
  return &Crypto{}
}

func (c *Crypto) Export(vm *goja.Runtime) goja.Value {
  c.vm = vm
  return vm.ToValue(c.createCryptoAPI())
}

func (c *Crypto) createCryptoAPI() map[string]interface{} {
  return map[string]interface{}{
    "createHash":      c.createHash,
    "createHmac":      c.createHmac,
    "timingSafeEqual": c.timingSafeEqual,
    "random":          c.random,
    "randomBytes":     c.random, // Alias for Node.js compatibility
    "uuid":            c.uuid,
  }
}

func (c *Crypto) createHash(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 1 {
    panic(c.vm.NewTypeError("createHash requires an algorithm argument"))
  }
    
  algorithm := call.Argument(0).String()
    
  obj := c.vm.NewObject()
    
  obj.Set("update", func(call goja.FunctionCall) goja.Value {
    if len(call.Arguments) < 1 {
      panic(c.vm.NewTypeError("update requires data argument"))
    }
    data := call.Argument(0).String()
    call.This.ToObject(c.vm).Set("_data", data)
    return call.This
  })
    
  obj.Set("digest", func(call goja.FunctionCall) goja.Value {
    encoding := "hex"
    if len(call.Arguments) > 0 {
      encoding = call.Argument(0).String()
    }
        
    thisObj := call.This.ToObject(c.vm)
    dataVal := thisObj.Get("_data")
    if dataVal == nil {
        return c.vm.ToValue("")
    }
        
    data := dataVal.String()
    var hashBytes []byte
        
    switch algorithm {
    case "md5":
      hash := md5.Sum([]byte(data))
      hashBytes = hash[:]
    case "sha1":
      hash := sha1.Sum([]byte(data))
      hashBytes = hash[:]
    case "sha256":
      hash := sha256.Sum256([]byte(data))
      hashBytes = hash[:]
    case "sha512":
      hash := sha512.Sum512([]byte(data))
      hashBytes = hash[:]
    default:
      panic(c.vm.NewTypeError(fmt.Sprintf("unsupported algorithm: %s", algorithm)))
    }
        
    switch encoding {
    case "hex":
      return c.vm.ToValue(hex.EncodeToString(hashBytes))
    case "base64":
      return c.vm.ToValue(base64.StdEncoding.EncodeToString(hashBytes))
    default:
      panic(c.vm.NewTypeError(fmt.Sprintf("unsupported encoding: %s", encoding)))
    }
  })
    
  return obj
}

func (c *Crypto) createHmac(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(c.vm.NewTypeError("createHmac requires algorithm and key arguments"))
  }
    
  algorithm := call.Argument(0).String()
  key := call.Argument(1).String()
    
  obj := c.vm.NewObject()
    
  obj.Set("update", func(call goja.FunctionCall) goja.Value {
    if len(call.Arguments) < 1 {
      panic(c.vm.NewTypeError("update requires data argument"))
    }
    data := call.Argument(0).String()
    call.This.ToObject(c.vm).Set("_data", data)
    return call.This
  })
    
  obj.Set("digest", func(call goja.FunctionCall) goja.Value {
    encoding := "hex"
    if len(call.Arguments) > 0 {
      encoding = call.Argument(0).String()
    }
        
    thisObj := call.This.ToObject(c.vm)
    dataVal := thisObj.Get("_data")
    if dataVal == nil {
      return c.vm.ToValue("")
    }
        
    data := dataVal.String()
    var h func() []byte
        
    switch algorithm {
    case "md5":
      h = func() []byte {
        mac := hmac.New(md5.New, []byte(key))
        mac.Write([]byte(data))
        return mac.Sum(nil)
      }
    case "sha1":
      h = func() []byte {
        mac := hmac.New(sha1.New, []byte(key))
        mac.Write([]byte(data))
        return mac.Sum(nil)
      }
    case "sha256":
      h = func() []byte {
        mac := hmac.New(sha256.New, []byte(key))
        mac.Write([]byte(data))
        return mac.Sum(nil)
      }
    case "sha512":
      h = func() []byte {
        mac := hmac.New(sha512.New, []byte(key))
        mac.Write([]byte(data))
        return mac.Sum(nil)
      }
    default:
      panic(c.vm.NewTypeError(fmt.Sprintf("unsupported algorithm: %s", algorithm)))
    }
        
    hashBytes := h()
        
    switch encoding {
    case "hex":
      return c.vm.ToValue(hex.EncodeToString(hashBytes))
    case "base64":
      return c.vm.ToValue(base64.StdEncoding.EncodeToString(hashBytes))
    default:
      panic(c.vm.NewTypeError(fmt.Sprintf("unsupported encoding: %s", encoding)))
    }
  })
    
  return obj
}

func (c *Crypto) timingSafeEqual(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(c.vm.NewTypeError("timingSafeEqual requires two arguments"))
  }
    
  a := call.Argument(0).String()
  b := call.Argument(1).String()
    
  // Convert strings to bytes
  aBytes := []byte(a)
  bBytes := []byte(b)
    
  // subtle.ConstantTimeCompare requires equal length
  if len(aBytes) != len(bBytes) {
    return c.vm.ToValue(false)
  }
    
  // Returns 1 if equal, 0 if not equal
  result := subtle.ConstantTimeCompare(aBytes, bBytes)
  return c.vm.ToValue(result == 1)
}

func (c *Crypto) uuid(call goja.FunctionCall) goja.Value {
  return c.vm.ToValue(uuid.New().String())
}

func (c *Crypto) random(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 1 {
    panic(c.vm.NewTypeError("random requires a size argument"))
  }

  size := int(call.Argument(0).ToInteger())
  if size < 0 {
    panic(c.vm.NewTypeError("size must be non-negative"))
  }
  if size > 65536 {
    panic(c.vm.NewTypeError("size must be less than 64kb"))
  }

  bytes := make([]byte, size)
  _, err := rand.Read(bytes)
  if err != nil {
    panic(c.vm.NewGoError(fmt.Errorf("failed to generate random bytes: %w", err)))
  }

  encoding := "hex"
  if len(call.Arguments) > 1 {
    encoding = call.Argument(1).String()
  }

  switch encoding {
  case "hex":
    return c.vm.ToValue(hex.EncodeToString(bytes))
  case "base64":
    return c.vm.ToValue(base64.StdEncoding.EncodeToString(bytes))
  case "raw":
    arr := c.vm.NewArray()
    for i, b := range bytes {
      arr.Set(fmt.Sprintf("%d", i), b)
    }
    return arr
  default:
    panic(c.vm.NewTypeError(fmt.Sprintf("unsupported encoding: %s (use 'hex', 'base64', or 'raw')", encoding)))
  }
}
