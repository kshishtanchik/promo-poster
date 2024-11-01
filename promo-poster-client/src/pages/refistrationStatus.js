//todo: отображение статуса регистрации на мероприятие
import React, {  } from "react";
import { Container } from "react-bootstrap";
import { useParams } from "react-router-dom";
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

function RegistrationStatus() {
    const { activityId } = useParams();
    
    return (
      <Container>
      <InputGroup>
        <InputGroup.Text>Имя</InputGroup.Text>
        <Form.Control as="textarea" id="name" aria-label="Имя" />
      </InputGroup>
      <InputGroup>
        <InputGroup.Text>Телефрон</InputGroup.Text>
        <Form.Control as="textarea" id="phone" type="phone" aria-label="Телефрон" />
      </InputGroup>
      <Form.Label htmlFor="telegram-url">Профиль в телеграм</Form.Label>
      <InputGroup className="mb-3">
        <InputGroup.Text id="telegram-alias">
        https://t.me/
        </InputGroup.Text>
        <Form.Control id="telegram-url" aria-describedby="telegram-alias" />
      </InputGroup>

      </Container>
    );
  }
  
  export default RegistrationStatus;