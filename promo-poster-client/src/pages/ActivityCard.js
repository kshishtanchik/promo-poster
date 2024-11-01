//todo: рредатирование мероприятия
import React from "react";
import Badge from "react-bootstrap/Badge";
import Button from "react-bootstrap/Button";
import { Container } from "react-bootstrap";
import Card from "react-bootstrap/Card";
import Stack from "react-bootstrap/Stack";
import { useParams } from "react-router-dom";
import Form from "react-bootstrap/Form";
import InputGroup from "react-bootstrap/InputGroup";

function NewActivityCard() {
  const { activityId } = useParams();
  let eventData = {
    activityType: "Тренировка",
    banner: "../../../images/blank_banner.jpg",
    address: "3я улица красивых молдавских партизан",
    startDate: new Date(),
  };

  function startEdit(e) {
    console.log(e.target.name);
    
  }

  function save(e) {
    console.log("save");
  }

  function selectPeriod(e) {
    console.log(e.target);
  }

  function loadBanner(e) {
    console.log(e.target)
    debugger
  }

  return (
    <Container>
      <Card className="text-center">
        <Card.Header>
          <InputGroup>
            <InputGroup.Text>Тип мероприятия</InputGroup.Text>
            <Form.Control type="text" as="input" id="name" aria-label="Имя" />
          </InputGroup>
        </Card.Header>
        {
          //todo: показывать картинку если есть
        }
        <Card.Img src={eventData.banner} onClick={startEdit} name="banner" />
        <Form.Group controlId="cover" className="mb-3">
          <Form.Label>Обложка</Form.Label>
          <Form.Control type="file" onChange={loadBanner} />
        </Form.Group>
        <Card.Body>
          <Card.Title>
            <InputGroup>
              <InputGroup.Text>Название</InputGroup.Text>
              <Form.Control
                type="text"
                as="input"
                id="title"
                aria-label="Название"
              />
            </InputGroup>
          </Card.Title>
          <InputGroup>
            <InputGroup.Text>Описание</InputGroup.Text>
            <Form.Control
              as="textarea"
              id="description"
              aria-label="Описание"
            />
          </InputGroup>
        </Card.Body>
        <Card.Footer className="text-muted">
          <InputGroup>
            <InputGroup.Text>Адресс</InputGroup.Text>
            <Form.Control
              as="input"
              id="address"
              aria-label="Адресс"
              defaultValue={eventData.address}
            />
          </InputGroup>
          <Stack direction="horizontal" gap={2}>
            <InputGroup className="p-2">
              <InputGroup.Text>Дата начала</InputGroup.Text>
              <Form.Control
                type="date"
                as="input"
                id="date"
                aria-label="Дата начала"
                defaultValue={eventData.startDate.toLocaleString("ru-RU")}
                onChange={selectPeriod}
              />
            </InputGroup>
            <div className="vr  ms-auto" />
            <InputGroup className="p-2 ms-auto">
              <InputGroup.Text>Во сколько:</InputGroup.Text>
              <Form.Control
                type="time"
                as="input"
                id="time"
                aria-label="Во сколько:"
                defaultValue={eventData.time}
              />
            </InputGroup>
          </Stack>
          <Form.Select
            aria-label="Default select example"
            onChange={selectPeriod}
          >
            <option>Периодичность</option>
            <option value="1">День</option>
            <option value="2">Неделя</option>
            <option value="3">Месяц</option>
            <option value="4">Год</option>
          </Form.Select>
        </Card.Footer>
        
        <Button variant="primary" onClick={save}>
            Сохранить
          </Button>
      </Card>
    </Container>
  );
}

export default NewActivityCard;
