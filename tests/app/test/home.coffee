supertest = require 'supertest'
expect = require('chai').expect

request = supertest 'http://ApiServer'
describe 'Homepage', ->
  it 'Should return standard greeting', (done) ->
    request
      .get '/'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.text).to.equal 'Hello, world!'
        done()
describe 'Ping endpoint', ->
  it 'Should respond to standard ping', (done) ->
    request
      .get '/ping'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.text).to.equal 'Pong'
        done()
describe 'Json reflection', ->
  it(
    'Should return identical Json in response as provided by request'
    (done) ->
      body = { greeting: 'Hello, world!' }
      request
        .post '/reflection'
        .send body
        .end (err, res) ->
          expect(err).to.equal null
          expect(res.body.greeting).to.equal body.greeting
          done()
  )
